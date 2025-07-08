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

package bkapi

import (
	"crypto/tls"
	"net/http"
	"time"

	"gopkg.in/h2non/gentleman.v2/plugins/auth"
	"gopkg.in/h2non/gentleman.v2/plugins/cookies"
	"gopkg.in/h2non/gentleman.v2/plugins/headers"
	"gopkg.in/h2non/gentleman.v2/plugins/proxy"
	"gopkg.in/h2non/gentleman.v2/plugins/query"
	"gopkg.in/h2non/gentleman.v2/plugins/redirect"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
	tlsplugin "gopkg.in/h2non/gentleman.v2/plugins/tls"
	"gopkg.in/h2non/gentleman.v2/plugins/transport"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal"
)

// OptTimeout defines the maximum amount of time a whole request process
// (including dial / request / redirect) can take.
func OptTimeout(duration time.Duration) define.BkApiOption {
	return internal.NewPluginOption(timeout.Request(duration))
}

// OptDialTimeout defines the maximum amount of time waiting for network dialing
func OptDialTimeout(duration, keepAlive time.Duration) define.BkApiOption {
	return internal.NewPluginOption(timeout.Dial(duration, keepAlive))
}

// OptTLShandshakeTimeout defines the maximum amount of time waiting for a TLS handshake
func OptTLShandshakeTimeout(duration time.Duration) define.BkApiOption {
	return internal.NewPluginOption(timeout.TLS(duration))
}

// OptBasicAuth defines an authorization basic header in the outgoing request
func OptBasicAuth(username, password string) define.BkApiOption {
	return internal.NewPluginOption(auth.Basic(username, password))
}

// OptBearerAuth defines an authorization bearer token header in the outgoing request
func OptBearerAuth(token string) define.BkApiOption {
	return internal.NewPluginOption(auth.Bearer(token))
}

// OptAddCookie adds a cookie to the request. Per RFC 6265 section 5.4, AddCookie does not
// attach more than one Cookie header field.
// That means all cookies, if any, are written into the same line, separated by semicolon.
func OptAddCookie(cookie *http.Cookie) define.BkApiOption {
	return internal.NewPluginOption(cookies.Add(cookie))
}

// OptDelAllCookies deletes all the cookies by deleting the Cookie header field.
func OptDelAllCookies() define.BkApiOption {
	return internal.NewPluginOption(cookies.DelAll())
}

// OptSetRequestHeader sets the header entries associated with key to the single element value.
// It replaces any existing values associated with key.
func OptSetRequestHeader(key string, value string) define.BkApiOption {
	return internal.NewPluginOption(headers.Set(key, value))
}

// OptDelRequestHeader deletes the header fields associated with key.
func OptDelRequestHeader(key string) define.BkApiOption {
	return internal.NewPluginOption(headers.Del(key))
}

// OptSetRequestHeaders sets the headers.
func OptSetRequestHeaders(header map[string]string) define.BkApiOption {
	return internal.NewPluginOption(headers.SetMap(header))
}

// OptProxies defines the proxy servers to be used based on the transport scheme
func OptProxies(servers map[string]string) define.BkApiOption {
	return internal.NewPluginOption(proxy.Set(servers))
}

// OptSetRequestQueryParam ets the query param key and value.
// It replaces any existing values.
func OptSetRequestQueryParam(key string, value string) define.BkApiOption {
	return internal.NewPluginOption(query.Set(key, value))
}

// OptAddRequestQueryParam adds the query param value to key.
// It appends to any existing values associated with key.
func OptAddRequestQueryParam(key string, value string) define.BkApiOption {
	return internal.NewPluginOption(query.Add(key, value))
}

// OptDelRequestQueryParam deletes the query param values associated with key.
func OptDelRequestQueryParam(key string) define.BkApiOption {
	return internal.NewPluginOption(query.Del(key))
}

// OptSetRequestQueryParams sets the query params.
func OptSetRequestQueryParams(params map[string]string) define.BkApiOption {
	return internal.NewPluginOption(query.SetMap(params))
}

// OptLimitRedirect defines in the maximum number of redirects that http.Client should follow.
func OptLimitRedirect(limit int) define.BkApiOption {
	return internal.NewPluginOption(redirect.Limit(limit))
}

// OptTransport sets a new HTTP transport for the outgoing request
func OptTransport(roundTripper http.RoundTripper) define.BkApiOption {
	return internal.NewPluginOption(transport.Set(roundTripper))
}

// OptTLS defines the request TLS connection config
func OptTLS(config *tls.Config) define.BkApiOption {
	return internal.NewPluginOption(tlsplugin.Config(config))
}
