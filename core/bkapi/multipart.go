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
	"gopkg.in/h2non/gentleman.v2/plugins/multipart"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal"
)

// MultipartFormFieldsBodyProvider provides request body as multipart form.
type MultipartFormFieldsBodyProvider struct {
	*FunctionalBodyProvider
}

// NewMultipartFormFieldsBodyProvider create a new MultipartFormFieldsBodyProvider
func NewMultipartFormFieldsBodyProvider() *MultipartFormFieldsBodyProvider {
	return &MultipartFormFieldsBodyProvider{
		FunctionalBodyProvider: NewFunctionalBodyProvider(func(operation define.Operation, v interface{}) error {
			values, ok := v.(map[string][]string)
			if !ok {
				return define.ErrorWrapf(define.ErrTypeNotMatch, "expected %T, but got %T", values, v)
			}

			fields := make(map[string]multipart.Values, len(values))
			for k, v := range values {
				fields[k] = multipart.Values(v)
			}

			operation.Apply(internal.NewPluginOption(multipart.Fields(fields)))

			return nil
		}),
	}
}

// MultipartFormBodyProvider provides request body as multipart form.
func MultipartFormBodyProvider() *MultipartFormFieldsBodyProvider {
	return NewMultipartFormFieldsBodyProvider()
}

// OptMultipartFormBodyProvider provides request body as multipart form.
func OptMultipartFormBodyProvider() *MultipartFormFieldsBodyProvider {
	return NewMultipartFormFieldsBodyProvider()
}
