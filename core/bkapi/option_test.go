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

package bkapi_test

import (
	"net/http"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/transport"
)

var _ = Describe("Option", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mock.MockRoundTripper
		operation    *internal.Operation
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mock.NewMockRoundTripper(ctrl)
		roundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Request:    req,
			}, nil
		}).AnyTimes()

		request := gentleman.NewRequest()
		request.Use(transport.Set(roundTripper))

		operation = internal.NewOperation("testing", request)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	applyOptions := func(opts ...define.OperationOption) *http.Response {
		response, err := operation.Apply(opts...).Request()
		Expect(err).To(BeNil())

		return response
	}

	It("OptAddRequestQueryParamList", func() {
		opt := bkapi.OptAddRequestQueryParamList("testing", []string{"a", "b", "c"})
		response := applyOptions(opt)

		Expect(response.Request.URL.Query().Encode()).To(Equal("testing=a&testing=b&testing=c"))
	})
})
