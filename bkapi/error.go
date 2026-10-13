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
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

// Error is returned when the API responds with a non-2xx status.
type Error struct {
	StatusCode int
	// Data contains the response's error value. Its fields and types are defined by the API.
	Data any
	// RequestID identifies the request in the logs of the gateway.
	RequestID string
}

func (e *Error) Error() string {
	return fmt.Sprintf("bkapi: status=%d error=%v request_id=%s", e.StatusCode, e.Data, e.RequestID)
}

// newError preserves the error value without assuming its structure. Older JSON responses
// are kept as a whole; non-JSON responses fall back to gateway headers and the HTTP status.
func newError(res *http.Response, body []byte) *Error {
	result := &Error{
		StatusCode: res.StatusCode,
		RequestID:  cmp.Or(res.Header.Get("X-Bkapi-Request-Id"), res.Header.Get("X-Request-Id")),
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err == nil && payload != nil {
		if data, ok := payload["error"]; ok {
			result.Data = data
		} else {
			result.Data = payload
		}
		return result
	}
	result.Data = map[string]any{
		"code":    res.Header.Get("X-Bkapi-Error-Code"),
		"message": cmp.Or(res.Header.Get("X-Bkapi-Error-Message"), http.StatusText(res.StatusCode)),
	}
	return result
}

// readBody reads the body of the response, and leaves it to be read again.
func readBody(res *http.Response) []byte {
	if res.Body == nil {
		return nil
	}
	body, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	res.Body = io.NopCloser(bytes.NewReader(body))
	return body
}

// errorBodyTransport reads the body of a non-2xx response into memory.
//
// gentleman runs the plugins of a client under a lock shared by all its requests, so checkStatus must not
// read a slow body itself. Reading it here keeps the request timeout and cancellation.
type errorBodyTransport struct {
	base http.RoundTripper
}

func (t *errorBodyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.base.RoundTrip(req)
	if err != nil || isSuccess(res.StatusCode) || res.Body == nil {
		return res, err
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return nil, err
	}
	res.Body = io.NopCloser(bytes.NewReader(body))
	return res, nil
}

func isSuccess(status int) bool {
	return status >= 200 && status < 300
}

// checkStatus returns the plugin that turns a non-2xx response into an *Error.
func checkStatus() plugin.Plugin {
	return plugin.NewResponsePlugin(func(ctx *context.Context, h context.Handler) {
		res := ctx.Response
		if isSuccess(res.StatusCode) {
			h.Next(ctx)
			return
		}

		h.Error(ctx, newError(res, readBody(res)))
	})
}
