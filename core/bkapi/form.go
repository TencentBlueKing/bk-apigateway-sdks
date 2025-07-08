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
	"net/url"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
)

// FormMarshalBodyProvider provides request body as urlencoded form.
type FormMarshalBodyProvider struct {
	*MarshalBodyProvider
}

// NewFormMarshalBodyProvider creates a new FormMarshalBodyProvider with marshal function.
func NewFormMarshalBodyProvider(marshaler func(v interface{}) ([]byte, error)) *FormMarshalBodyProvider {
	return &FormMarshalBodyProvider{
		MarshalBodyProvider: NewMarshalBodyProvider("application/x-www-form-urlencoded", marshaler),
	}
}

// FormBodyProvider is a function to set form body from map[string][]string
func FormBodyProvider() *FormMarshalBodyProvider {
	return NewFormMarshalBodyProvider(func(v interface{}) ([]byte, error) {
		values, ok := v.(map[string][]string)
		if !ok {
			return nil, define.ErrorWrapf(define.ErrTypeNotMatch, "expected %T, but got %T", values, v)
		}

		forms := url.Values(values)
		return []byte(forms.Encode()), nil
	})
}

// OptFormBodyProvider is a function to set form body
func OptFormBodyProvider() *FormMarshalBodyProvider {
	return FormBodyProvider()
}
