/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2017-2021 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package manager_test

import (
	"net/http"
	"time"

	manager "github.com/TencentBlueKing/bk-apigateway-sdks/apigw-manager"
	apigateway "github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-bk-apigateway"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gock.v1"
)

var _ = Describe("Publickey", func() {
	var (
		config   bkapi.ClientConfig
		provider *manager.PublicKeyMemoryCache
	)

	BeforeEach(func() {
		config = bkapi.ClientConfig{
			Endpoint: "http://example.com",
		}
		provider = manager.NewPublicKeyMemoryCache(
			config, time.Hour,
			func(apiName string, config bkapi.ClientConfig) (*manager.Manager, error) {
				return manager.NewManager(
					apiName,
					config,
					nil,
					func(configProvider define.ClientConfigProvider, opts ...define.BkApiClientOption) (*apigateway.Client, error) {
						opts = append(opts, bkapi.OptTransport(gock.NewTransport()))
						return apigateway.New(configProvider, opts...)
					},
				)
			},
		)
	})

	It("should cache the public key", func() {
		count := 0
		gock.New(config.Endpoint).
			Get("/api/v1/apis/testing/public_key/").
			AddMatcher(func(_ *http.Request, _ *gock.Request) (bool, error) {
				count++
				return true, nil
			}).
			Reply(200).
			JSON(map[string]interface{}{
				"code": 0,
				"data": map[string]interface{}{
					"public_key": "public_key",
				},
			})
		defer gock.Off()

		for i := 0; i < 10; i++ {
			publicKey, err := provider.ProvidePublicKey("testing")
			Expect(err).To(BeNil())
			Expect(publicKey).To(Equal("public_key"))
		}

		Expect(count).To(Equal(1))
	})
})
