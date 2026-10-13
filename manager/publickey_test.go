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

package manager_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gock "gopkg.in/h2non/gock.v1"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/bkapi"
	manager "github.com/TencentBlueKing/bk-apigateway-sdks/v2/manager"
)

var _ = Describe("Publickey", func() {
	var (
		config   bkapi.Config
		provider *manager.PublicKeyMemoryCache
	)

	BeforeEach(func() {
		config = bkapi.Config{Endpoint: "http://example.com", Transport: gock.NewTransport()}
		provider = manager.NewPublicKeyMemoryCache(
			config, time.Hour,
			func(gatewayName string, config bkapi.Config) (*manager.Manager, error) {
				return manager.NewManager(
					gatewayName,
					config,
					nil,
				)
			},
		)
	})

	It("should cache the public key", func() {
		count := 0
		gock.New(config.Endpoint).
			Get("/api/v2/sync/gateways/testing/public_key/").
			AddMatcher(func(_ *http.Request, _ *gock.Request) (bool, error) {
				count++
				return true, nil
			}).
			Reply(200).
			JSON(map[string]any{
				"data": map[string]any{
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
