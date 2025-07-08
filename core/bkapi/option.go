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
	"net/http"

	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal"
)

// OptAddRequestQueryParamList adds the query param value list to key.
// It appends to any existing values associated with key.
func OptAddRequestQueryParamList(key string, values []string) define.BkApiOption {
	return internal.NewPluginOption(plugin.NewRequestPlugin(func(ctx *context.Context, h context.Handler) {
		query := ctx.Request.URL.Query()

		for _, value := range values {
			query.Add(key, value)
		}

		ctx.Request.URL.RawQuery = query.Encode()
		h.Next(ctx)
	}))
}

// OptRequestCallback sets the callback function for the request.
func OptRequestCallback(fn func(request *http.Request) *http.Request) define.BkApiOption {
	return internal.NewPluginOption(plugin.NewRequestPlugin(func(ctx *context.Context, h context.Handler) {
		ctx.Request = fn(ctx.Request)
		h.Next(ctx)
	}))
}

// OptResponseCallback sets the callback function for the response.
func OptResponseCallback(fn func(response *http.Response) *http.Response) define.BkApiOption {
	return internal.NewPluginOption(plugin.NewResponsePlugin(func(ctx *context.Context, h context.Handler) {
		ctx.Response = fn(ctx.Response)
		h.Next(ctx)
	}))
}

// OptErrorCallback sets the callback function for the error.
func OptErrorCallback(fn func(err error) error) define.BkApiOption {
	return internal.NewPluginOption(plugin.NewErrorPlugin(func(ctx *context.Context, h context.Handler) {
		ctx.Error = fn(ctx.Error)
		h.Next(ctx)
	}))
}
