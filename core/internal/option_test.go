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

package internal_test

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal/mock"
	"github.com/TencentBlueKing/gopkg/logging"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gentleman.v2"
)

var _ = Describe("Option", func() {
	Context("PluginOption", func() {
		var (
			ctrl             *gomock.Controller
			option           *internal.PluginOption
			pluginA, pluginB *mock.MockPlugin
			client           *gentleman.Client
			bkapiClient      *mock.MockBkApiClient
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			pluginA = mock.NewMockPlugin(ctrl)
			pluginB = mock.NewMockPlugin(ctrl)
			option = internal.NewPluginOption(pluginA, pluginB)
			client = gentleman.New()
			bkapiClient = mock.NewMockBkApiClient(ctrl)
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should apply to operation", func() {
			request := client.Request()
			operation := internal.NewOperation("", bkapiClient, request)
			operation.Apply(option)
			Expect(operation.GetError()).To(BeNil())

			stacks := request.Middleware.GetStack()
			Expect(stacks).To(ContainElement(pluginA))
			Expect(stacks).To(ContainElement(pluginB))
		})

		It("should apply to client", func() {
			clientConfig := mock.NewMockClientConfig(ctrl)
			clientConfig.EXPECT().GetUrl().Return("").AnyTimes()
			clientConfig.EXPECT().GetAuthorizationHeaders().Return(nil).AnyTimes()
			clientConfig.EXPECT().GetLogger().Return(logging.GetLogger("")).AnyTimes()

			cli := internal.NewBkApiClient("", client, nil, clientConfig)

			Expect(cli.Apply(option)).To(Succeed())

			stacks := client.Middleware.GetStack()
			Expect(stacks).To(ContainElement(pluginA))
			Expect(stacks).To(ContainElement(pluginB))
		})
	})
})
