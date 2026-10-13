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
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/bkapi"
)

var _ = Describe("ConfigFromEnv", func() {
	BeforeEach(func() {
		// unset the variables, which are restored after the spec
		for _, key := range []string{
			"BK_API_URL_TMPL", "BK_APP_CODE", "BKPAAS_APP_ID", "APP_CODE", "BK_APP_SECRET", "BKPAAS_APP_SECRET",
			"SECRET_KEY", "BKPAAS_APP_TENANT_ID", "BK_APP_TENANT_ID",
		} {
			GinkgoT().Setenv(key, "")
			Expect(os.Unsetenv(key)).To(Succeed())
		}
	})

	DescribeTable("should read the config from the environment variables",
		func(env map[string]string, want bkapi.Config) {
			for key, value := range env {
				GinkgoT().Setenv(key, value)
			}
			Expect(bkapi.ConfigFromEnv("my-gateway", "stag")).To(Equal(want))
		},
		Entry("of an app deployed on paas",
			map[string]string{
				"BK_API_URL_TMPL":      "http://bkapi.example.com/api/{api_name}/",
				"BKPAAS_APP_ID":        "paas-app",
				"BKPAAS_APP_SECRET":    "paas-secret",
				"BKPAAS_APP_TENANT_ID": "paas-tenant",
				"BK_APP_TENANT_ID":     "app-tenant",
			},
			bkapi.Config{
				Endpoint:  "http://bkapi.example.com/api/my-gateway/stag",
				AppCode:   "paas-app",
				AppSecret: "paas-secret",
				TenantID:  "paas-tenant",
			}),
		Entry("of a global tenant app deployed on paas",
			map[string]string{"BKPAAS_APP_TENANT_ID": "", "BK_APP_TENANT_ID": "app-tenant"},
			bkapi.Config{TenantID: "system"}),
		Entry("with the gateway name in the url template",
			map[string]string{"BK_API_URL_TMPL": "http://{gateway_name}.apigw.example.com"},
			bkapi.Config{Endpoint: "http://my-gateway.apigw.example.com/stag"}),
		Entry("of an app not deployed on paas",
			map[string]string{
				"BK_APP_CODE":      "app",
				"BKPAAS_APP_ID":    "paas-app",
				"BK_APP_SECRET":    "secret",
				"SECRET_KEY":       "legacy-secret",
				"BK_APP_TENANT_ID": "app-tenant",
			},
			bkapi.Config{AppCode: "app", AppSecret: "secret", TenantID: "app-tenant"}),
	)
})
