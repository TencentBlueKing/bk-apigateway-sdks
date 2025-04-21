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
	"fmt"
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

var _ = Describe("Provider", func() {
	var ctrl *gomock.Controller

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("MarshalBodyProvider", func() {
		var operation *mock.MockOperation

		BeforeEach(func() {
			operation = mock.NewMockOperation(ctrl)
		})

		It("should marshal the data", func() {
			result := "hello world"

			operation.EXPECT().SetContentType("text/plain").Return(operation)
			operation.EXPECT().SetContentLength(int64(len(result))).Return(operation)
			operation.EXPECT().SetBodyReader(gomock.Any()).DoAndReturn(func(reader io.Reader) define.Operation {
				body, err := ioutil.ReadAll(reader)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal(result))

				return operation
			})

			provider := bkapi.NewMarshalBodyProvider("text/plain", func(v interface{}) ([]byte, error) {
				return []byte(fmt.Sprintf("hello %s", v)), nil
			})
			Expect(provider.ProvideBody(operation, "world")).To(BeNil())
		})

		It("should set an empty body", func() {
			provider := bkapi.NewMarshalBodyProvider("text/plain", func(v interface{}) ([]byte, error) {
				panic("should not be called")
			})

			Expect(provider.ProvideBody(operation, nil)).To(BeNil())
		})
	})

	Context("UnmarshalResultProvider", func() {
		It("should unmarshal the result", func() {
			var reader io.Reader

			provider := bkapi.NewUnmarshalResultProvider(func(body io.Reader, v interface{}) error {
				*v.(*io.Reader) = body

				return nil
			})

			Expect(provider.ProvideResult(&http.Response{
				Body: ioutil.NopCloser(strings.NewReader("hello world")),
			}, &reader)).To(BeNil())

			result, err := ioutil.ReadAll(reader)
			Expect(err).To(BeNil())
			Expect(string(result)).To(Equal("hello world"))
		})

		It("should not unmarshal the result", func() {
			provider := bkapi.NewUnmarshalResultProvider(func(body io.Reader, v interface{}) error {
				panic("should not be called")
			})

			Expect(provider.ProvideResult(&http.Response{
				Body: ioutil.NopCloser(strings.NewReader("hello world")),
			}, nil)).To(BeNil())
		})
	})

	DescribeTable("provider as BkApiClientOption", func(provider interface{}) {
		opt, ok := provider.(define.BkApiClientOption)
		Expect(ok).To(BeTrue())

		client := mock.NewMockBkApiClient(ctrl)
		client.EXPECT().AddOperationOptions(opt).Return(nil)

		Expect(opt.ApplyToClient(client)).To(BeNil())
	},
		Entry("MarshalBodyProvider", bkapi.NewMarshalBodyProvider("", nil)),
		Entry("UnmarshalResultProvider", bkapi.NewUnmarshalResultProvider(nil)),
		Entry("FunctionalBodyProvider", bkapi.NewFunctionalBodyProvider(nil)),
	)

	DescribeTable("BodyProvider as OperationOption", func(provider define.BodyProvider) {
		opt, ok := provider.(define.OperationOption)
		Expect(ok).To(BeTrue())

		operation := mock.NewMockOperation(ctrl)
		operation.EXPECT().SetBodyProvider(provider).Return(nil)

		Expect(opt.ApplyToOperation(operation)).To(BeNil())
	},
		Entry("MarshalBodyProvider", bkapi.NewMarshalBodyProvider("", nil)),
		Entry("FunctionalBodyProvider", bkapi.NewFunctionalBodyProvider(nil)),
	)

	DescribeTable("ResultProvider as OperationOption", func(provider define.ResultProvider) {
		opt, ok := provider.(define.OperationOption)
		Expect(ok).To(BeTrue())

		operation := mock.NewMockOperation(ctrl)
		operation.EXPECT().SetResultProvider(provider).Return(nil)

		Expect(opt.ApplyToOperation(operation)).To(BeNil())
	},
		Entry("UnmarshalResultProvider", bkapi.NewUnmarshalResultProvider(nil)),
	)
})
