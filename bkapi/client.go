/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

// Package bkapi creates the gentleman clients to call the APIs through BlueKing API Gateway.
//
// A client handles what all the APIs have in common with gentleman plugins: authentication, tenant,
// timeout, response status and logging. Requests and responses are built and read with gentleman:
//
//	client, err := bkapi.New(bkapi.ConfigFromEnv("my-gateway", "prod"))
//	res, err := client.Get().AddPath("/users/:id/").Param("id", id).Use(bkapi.WithContext(ctx)).Send()
//	err = res.JSON(&user)
//
// Paths must be added with AddPath, as Path replaces the path of the endpoint.
package bkapi

import (
	"cmp"
	stdcontext "context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	gentleman "gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
	"gopkg.in/h2non/gentleman.v2/plugins/headers"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
	"gopkg.in/h2non/gentleman.v2/plugins/transport"
)

const defaultTimeout = 60 * time.Second

// Config is the configuration of a client.
type Config struct {
	// Endpoint is the URL of a gateway stage, such as https://bkapi.example.com/api/my-gateway/prod.
	Endpoint string

	// AppCode, AppSecret and AccessToken are sent in the X-Bkapi-Authorization header.
	AppCode     string
	AppSecret   string
	AccessToken string
	// TenantID is sent in the X-Bk-Tenant-Id header.
	TenantID string

	// Timeout limits the time of each request, including reading the response. Defaults to 60s.
	Timeout time.Duration
	// Transport sends the requests. Defaults to gentleman.DefaultTransport.
	Transport http.RoundTripper
	// Logger logs the requests: successful ones at debug level, failed ones at warn or error level.
	// Defaults to slog.Default().
	Logger *slog.Logger
}

// New creates a gentleman client to call the APIs of the gateway stage. It is safe for concurrent use.
// A non-2xx response is returned as an *Error by Send.
func New(config Config) (*gentleman.Client, error) {
	endpoint, err := url.Parse(config.Endpoint)
	if err != nil || endpoint.Scheme == "" || endpoint.Host == "" {
		return nil, fmt.Errorf("bkapi: invalid endpoint %q", config.Endpoint)
	}

	authorization, err := json.Marshal(struct {
		AppCode     string `json:"bk_app_code,omitempty"`
		AppSecret   string `json:"bk_app_secret,omitempty"`
		AccessToken string `json:"access_token,omitempty"`
	}{config.AppCode, config.AppSecret, config.AccessToken})
	if err != nil {
		return nil, err
	}

	base := cmp.Or[http.RoundTripper](config.Transport, gentleman.DefaultTransport)

	client := gentleman.New().
		URL(strings.TrimSuffix(config.Endpoint, "/")).
		Use(headers.Set("X-Bkapi-Authorization", string(authorization))).
		Use(timeout.Request(cmp.Or(config.Timeout, defaultTimeout))).
		Use(transport.Set(&errorBodyTransport{base: base})).
		Use(logging(cmp.Or(config.Logger, slog.Default()))).
		Use(checkStatus())
	if config.TenantID != "" {
		client.Use(headers.Set("X-Bk-Tenant-Id", config.TenantID))
	}
	return client, nil
}

// WithContext returns the plugin that cancels the request when ctx is done.
func WithContext(ctx stdcontext.Context) plugin.Plugin {
	return plugin.NewRequestPlugin(func(c *context.Context, h context.Handler) {
		h.Next(c.SetCancelContext(ctx))
	})
}

const startKey = "bkapi.start"

// logging returns the plugin that logs the responses, at debug level for 2xx, warn for 4xx and error for others.
// It must run before checkStatus, which stops the response phase of a non-2xx response.
//
// The errors without a response, such as network errors and timeouts, are only returned to the caller, who knows
// the context to log them. The headers, query and bodies are not logged, as they may contain secrets.
func logging(logger *slog.Logger) plugin.Plugin {
	p := plugin.New()
	p.SetHandlers(plugin.Handlers{
		"request": func(ctx *context.Context, h context.Handler) {
			ctx.Set(startKey, time.Now())
			h.Next(ctx)
		},
		"response": func(ctx *context.Context, h context.Handler) {
			logResponse(logger, ctx)
			h.Next(ctx)
		},
	})
	return p
}

func logResponse(logger *slog.Logger, ctx *context.Context) {
	start, _ := ctx.Get(startKey).(time.Time)
	res := ctx.Response
	attrs := []slog.Attr{
		slog.String("method", ctx.Request.Method),
		slog.String("path", ctx.Request.URL.Path),
		slog.Int("status", res.StatusCode),
		slog.String("request_id", cmp.Or(res.Header.Get("X-Bkapi-Request-Id"), res.Header.Get("X-Request-Id"))),
		slog.Duration("duration", time.Since(start)),
	}

	level := slog.LevelDebug
	if !isSuccess(res.StatusCode) {
		level = slog.LevelError
		if res.StatusCode >= 400 && res.StatusCode < 500 {
			level = slog.LevelWarn
		}
		// the body is already in memory, see errorBodyTransport
		apiErr := newError(res, readBody(res))
		attrs = append(attrs, slog.Any("error", apiErr.Data))
	}
	if traceparent := ctx.Request.Header.Get("Traceparent"); traceparent != "" {
		attrs = append(attrs, slog.String("traceparent", traceparent))
	}
	logger.LogAttrs(ctx.Request.Context(), level, "bkapi request", attrs...)
}
