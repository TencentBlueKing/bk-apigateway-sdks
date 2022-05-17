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

package demo_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/demo"
)

type mockTransport struct {
	response *http.Response
}

func (t mockTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return t.response, nil
}

func newMockTransport() *mockTransport {
	return &mockTransport{
		response: &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Header:     http.Header{},
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"demo": true}`)),
		},
	}
}

func Benchmark_Demo_Request(b *testing.B) {
	client, _ := demo.New(bkapi.ClientConfig{}, bkapi.OptTransport(newMockTransport()))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var result map[string]interface{}
			_, _ = client.Anything().
				SetResultProvider(bkapi.JsonResultProvider()).
				SetResult(&result).
				Request()
		}
	})
}
