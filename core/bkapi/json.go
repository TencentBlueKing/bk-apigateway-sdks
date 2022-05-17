/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2017 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package bkapi

import (
	"encoding/json"
	"io"
)

// JsonMarshalBodyProvider provides request body as json.
type JsonMarshalBodyProvider struct {
	*MarshalBodyProvider
}

// NewJsonMarshalBodyProvider creates a new JsonMarshalBodyProvider with marshal function.
func NewJsonMarshalBodyProvider(marshaler func(v interface{}) ([]byte, error)) *JsonMarshalBodyProvider {
	return &JsonMarshalBodyProvider{
		MarshalBodyProvider: NewMarshalBodyProvider("application/json", marshaler),
	}
}

// JsonBodyProvider creates a new JsonMarshalBodyProvider with default marshal function.
func JsonBodyProvider() *JsonMarshalBodyProvider {
	return NewJsonMarshalBodyProvider(json.Marshal)
}

// OptJsonBodyProvider is a option for json body provider.
func OptJsonBodyProvider() *JsonMarshalBodyProvider {
	return JsonBodyProvider()
}

// JsonUnmarshalResultProvider provides result from json.
type JsonUnmarshalResultProvider struct {
	*UnmarshalResultProvider
}

// NewJsonUnmarshalResultProvider creates a new JsonUnmarshalResultProvider with unmarshal function.
func NewJsonUnmarshalResultProvider(
	unmarshaler func(body io.Reader, v interface{}) error,
) *JsonUnmarshalResultProvider {
	return &JsonUnmarshalResultProvider{
		UnmarshalResultProvider: NewUnmarshalResultProvider(unmarshaler),
	}
}

// JsonResultProvider creates a new JsonUnmarshalResultProvider with default unmarshal function.
func JsonResultProvider() *JsonUnmarshalResultProvider {
	return NewJsonUnmarshalResultProvider(func(body io.Reader, v interface{}) error {
		return json.NewDecoder(body).Decode(v)
	})
}

// OptJsonResultProvider is a option for json result provider.
func OptJsonResultProvider() *JsonUnmarshalResultProvider {
	return JsonResultProvider()
}
