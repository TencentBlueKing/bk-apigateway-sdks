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

package apigateway_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/TencentBlueKing/bk-apigateway-sdks/apigateway"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
)

var _ = Describe("Client", func() {
	It("should create a client by config", func() {
		client, err := apigateway.New(bkapi.ClientConfig{
			BkApiUrlTmpl: "https://{api_name}.example.com/",
			Stage:        "prod",
			AccessToken:  "access_token",
			AppCode:      "app_code",
			AppSecret:    "app_secret",
		})
		Expect(err).To(BeNil())

		Expect(client.Name()).To(Equal("bk-apigateway"))
	})
})
