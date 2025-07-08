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
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal/mock"
)

var _ = Describe("Operation", func() {
	Context("OperationConfig", func() {
		It("should clone a config", func() {
			config := bkapi.OperationConfig{
				Name:   "test",
				Path:   "/test",
				Method: "GET",
			}

			providedConfig := config.ProvideConfig()

			config.Name = ""
			config.Path = ""
			config.Method = ""

			Expect(providedConfig.GetName()).To(Equal("test"))
			Expect(providedConfig.GetPath()).To(Equal("/test"))
			Expect(providedConfig.GetMethod()).To(Equal("GET"))
		})
	})

	Context("OperationOption", func() {
		var ctrl *gomock.Controller

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should apply to client", func() {
			option := bkapi.NewOperationOption(nil)

			client := mock.NewMockBkApiClient(ctrl)
			client.EXPECT().AddOperationOptions(option).Return(nil)

			Expect(option.ApplyToClient(client)).To(Succeed())
		})

		It("should apply to operation", func() {
			operation := mock.NewMockOperation(ctrl)

			called := false
			option := bkapi.NewOperationOption(func(op define.Operation) error {
				called = true
				Expect(op).To(Equal(operation))

				return nil
			})

			Expect(option.ApplyToOperation(operation)).To(Succeed())
			Expect(called).To(BeTrue())
		})
	})

	Context("Option", func() {
		var (
			ctrl      *gomock.Controller
			operation *mock.MockOperation
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())

			operation = mock.NewMockOperation(ctrl)
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should set the body", func() {
			body := make(map[string]interface{})
			operation.EXPECT().SetBody(body).Return(nil)

			option := bkapi.OptSetRequestBody(body)
			Expect(option.ApplyToOperation(operation)).To(Succeed())
		})

		It("should set the result", func() {
			result := make(map[string]interface{})
			operation.EXPECT().SetResult(result).Return(nil)

			option := bkapi.OptSetRequestResult(result)
			Expect(option.ApplyToOperation(operation)).To(Succeed())
		})
	})
})
