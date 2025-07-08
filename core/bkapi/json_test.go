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

package bkapi_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal/mock"
)

var _ = Describe("Json", func() {
	var ctrl *gomock.Controller

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should provide urlencoded form", func() {
		result := `{"hello":"world"}`
		operation := mock.NewMockOperation(ctrl)
		operation.EXPECT().SetContentType("application/json").Return(operation)
		operation.EXPECT().SetContentLength(int64(len(result))).Return(operation)
		operation.EXPECT().SetBodyReader(gomock.Any()).DoAndReturn(func(body io.Reader) define.Operation {
			data, err := ioutil.ReadAll(body)
			Expect(err).To(BeNil())
			Expect(string(data)).To(Equal(result))
			return operation
		})

		provider := bkapi.JsonBodyProvider()
		Expect(provider.ProvideBody(operation, map[string]interface{}{
			"hello": "world",
		})).To(BeNil())
	})

	It("should decode json result", func() {
		var result map[string]interface{}
		provider := bkapi.JsonResultProvider()
		Expect(provider.ProvideResult(&http.Response{
			Body: ioutil.NopCloser(strings.NewReader(`{"hello":"world"}`)),
		}, &result)).To(BeNil())

		Expect(result["hello"]).To(Equal("world"))
	})
})
