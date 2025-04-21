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
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal/mock"
)

var _ = Describe("Config", func() {
	Context("ClientConfigRegistry", func() {
		var (
			defaultConfig ClientConfig
			registry      *ClientConfigRegistry
		)

		BeforeEach(func() {
			defaultConfig = ClientConfig{
				BkApiUrlTmpl: "http://{api_name}.example.com/",
			}
			registry = NewClientConfigRegistry()
			Expect(registry.RegisterDefaultConfig(defaultConfig)).To(Succeed())
		})

		It("should return default config when it is not found", func() {
			apiName := "do-not-exist"
			config := registry.ProvideConfig(apiName)

			Expect(config.GetName()).To(Equal(apiName))
			Expect(config.GetUrl()).To(Equal("http://do-not-exist.example.com/prod/"))
		})

		It("should return a registered config", func() {
			apiName := "should-exist"
			clientConfig := ClientConfig{
				Endpoint: "http://special.example.com/",
			}

			Expect(registry.RegisterClientConfig(apiName, clientConfig)).To(Succeed())
			config := registry.ProvideConfig(apiName)

			Expect(config.GetName()).To(Equal(apiName))
			Expect(config.GetUrl()).To(Equal(clientConfig.Endpoint))
		})

		It("support to register by a custom config provider", func() {
			ctrl := gomock.NewController(GinkgoT())
			defer ctrl.Finish()

			provider := mock.NewMockClientConfigProvider(ctrl)
			providedConfig := mock.NewMockClientConfig(ctrl)
			provider.EXPECT().ProvideConfig(gomock.Any()).Return(providedConfig)

			apiName := "my-config"
			Expect(registry.RegisterClientConfig(apiName, provider)).To(Succeed())
			config := registry.ProvideConfig(apiName)

			Expect(config).To(Equal(providedConfig))
		})
	})
})
