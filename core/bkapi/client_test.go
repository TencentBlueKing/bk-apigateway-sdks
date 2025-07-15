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
	"net/http"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal/mock"
)

var _ = Describe("Client", func() {
	var ctrl *gomock.Controller

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("BkApiClient", func() {
		var (
			apiName              = "testing"
			configProvider       *mock.MockClientConfigProvider
			config               *mock.MockClientConfig
			roundTripper         *mock.MockRoundTripper
			roundTripperOpt      define.BkApiOption
			url                  string
			authorizationHeaders map[string]string
		)

		BeforeEach(func() {
			url = "http://api.example.com/"
			configProvider = mock.NewMockClientConfigProvider(ctrl)
			config = mock.NewMockClientConfig(ctrl)
			roundTripper = mock.NewMockRoundTripper(ctrl)
			roundTripperOpt = bkapi.OptTransport(roundTripper)
			authorizationHeaders = make(map[string]string)

			configProvider.EXPECT().ProvideConfig(apiName).Return(config).AnyTimes()
			config.EXPECT().GetUrl().DoAndReturn(func() string {
				return url
			}).AnyTimes()
			config.EXPECT().GetAuthorizationHeaders().DoAndReturn(func() map[string]string {
				return authorizationHeaders
			}).AnyTimes()
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		getMockRequest := func(client define.BkApiClient) (request *http.Request) {
			roundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
				request = req
				return &http.Response{}, nil
			})

			operationConfig := mock.NewMockOperationConfig(ctrl)
			operationConfig.EXPECT().GetName().Return("").AnyTimes()
			operationConfig.EXPECT().GetMethod().Return("").AnyTimes()
			operationConfig.EXPECT().GetPath().Return("").AnyTimes()

			operationConfigProvider := mock.NewMockOperationConfigProvider(ctrl)
			operationConfigProvider.EXPECT().ProvideConfig().Return(operationConfig).AnyTimes()

			operation := client.NewOperation(operationConfigProvider)
			_, err := operation.Request()
			Expect(err).To(BeNil())

			return request
		}

		It("should apply option by argument", func() {
			config.EXPECT().GetLogger().Return(nil).AnyTimes()
			config.EXPECT().GetClientOptions().Return(nil)

			client, err := bkapi.NewBkApiClient(apiName, configProvider, roundTripperOpt)
			Expect(err).To(BeNil())

			request := getMockRequest(client)
			Expect(request).NotTo(BeNil())
		})

		It("should apply option by config", func() {
			config.EXPECT().GetLogger().Return(nil).AnyTimes()
			config.EXPECT().GetClientOptions().Return([]define.BkApiClientOption{roundTripperOpt})

			client, err := bkapi.NewBkApiClient(apiName, configProvider)
			Expect(err).To(BeNil())

			request := getMockRequest(client)
			Expect(request).NotTo(BeNil())
		})
	})

	Context("ClientConfig", func() {
		It("should clone a new config", func() {
			config := bkapi.ClientConfig{}
			Expect(config.GetName()).To(Equal(""))

			providedConfig := config.ProvideConfig("testing").(*bkapi.ClientConfig)
			Expect(providedConfig.GetName()).To(Equal("testing"))

			Expect(config.GetName()).To(Equal(""))
		})

		It("should return endpoint as url", func() {
			config := bkapi.ClientConfig{
				Endpoint: "http://example.com",
			}

			Expect(config.ProvideConfig("testing").GetUrl()).To(Equal("http://example.com/"))
		})

		It("should render endpoint with params", func() {
			config := bkapi.ClientConfig{
				BkApiUrlTmpl: "http://{api_name}.example.com/",
				Stage:        "prod",
			}

			Expect(config.ProvideConfig("testing").GetUrl()).To(Equal("http://testing.example.com/prod/"))
		})

		It("should not return authorization headers when related params are empty", func() {
			config := bkapi.ClientConfig{
				AccessToken:         "",
				AuthorizationJWT:    "",
				AppCode:             "",
				AppSecret:           "",
				AuthorizationParams: nil,
			}

			Expect(config.GetAuthorizationHeaders()).To(BeEmpty())
		})

		It("should return access token authorization headers", func() {
			config := bkapi.ClientConfig{
				AccessToken:      "access_token",
				AuthorizationJWT: "jwt",
			}

			Expect(config.GetAuthorizationHeaders()).To(Equal(map[string]string{
				"X-Bkapi-Authorization": `{"access_token":"access_token","jwt":"jwt"}`,
				"X-Bk-Tenant-Id":        "",
			}))
		})

		It("should return app code authorization headers", func() {
			config := bkapi.ClientConfig{
				AppCode:   "app_code",
				AppSecret: "app_secret",
			}

			Expect(config.GetAuthorizationHeaders()).To(Equal(map[string]string{
				"X-Bkapi-Authorization": `{"bk_app_code":"app_code","bk_app_secret":"app_secret"}`,
				"X-Bk-Tenant-Id":        "",
			}))
		})

		It("should return common authorization headers", func() {
			config := bkapi.ClientConfig{
				AuthorizationParams: map[string]string{
					"bk_token": "token",
				},
			}

			Expect(config.GetAuthorizationHeaders()).To(Equal(map[string]string{
				"X-Bkapi-Authorization": `{"bk_token":"token"}`,
				"X-Bk-Tenant-Id":        "",
			}))
		})

		It("should return authorization headers marshal by custom marshaler", func() {
			config := bkapi.ClientConfig{
				AccessToken: "access_token",
				JsonMarshaler: func(v interface{}) ([]byte, error) {
					return []byte(`{"access_token": "access_token"}`), nil
				},
			}

			Expect(config.GetAuthorizationHeaders()).To(Equal(map[string]string{
				"X-Bkapi-Authorization": `{"access_token": "access_token"}`,
				"X-Bk-Tenant-Id":        "",
			}))
		})

		It("should set stage by default", func() {
			config := bkapi.ClientConfig{}
			providedConfig := config.ProvideConfig("testing").(*bkapi.ClientConfig)

			Expect(providedConfig.Stage).To(Equal("prod"))
		})

		It("should use endpoint directly when it is not empty", func() {
			config := bkapi.ClientConfig{
				Endpoint: "http://api.example.com/",
			}
			providedConfig := config.ProvideConfig("testing").(*bkapi.ClientConfig)

			Expect(providedConfig.Endpoint).To(Equal("http://api.example.com/"))
			Expect(providedConfig.Stage).To(Equal(""))
			Expect(providedConfig.BkApiUrlTmpl).To(Equal(""))
		})

		It("should set endpoint by BkApiUrlTmpl", func() {
			config := bkapi.ClientConfig{
				Stage:        "test",
				BkApiUrlTmpl: "http://{api_name}.example.com/",
			}
			providedConfig := config.ProvideConfig("testing").(*bkapi.ClientConfig)

			Expect(providedConfig.Endpoint).To(Equal("http://testing.example.com/test"))
		})

		It("should set endpoint by env BK_API_URL_TMPL", func() {
			config := bkapi.ClientConfig{
				Stage: "test",
				Getenv: func(k string) string {
					if k == "BK_API_URL_TMPL" {
						return "http://{api_name}.example.com/"
					}
					return ""
				},
			}
			providedConfig := config.ProvideConfig("testing").(*bkapi.ClientConfig)

			Expect(providedConfig.Endpoint).To(Equal("http://testing.example.com/test"))
		})

		It("should set endpoint by env BK_API_STAGE_URL_TMPL", func() {
			config := bkapi.ClientConfig{
				Stage: "dev",
				Getenv: func(k string) string {
					if k == "BK_API_STAGE_URL_TMPL" {
						return "http://{stage}-{api_name}.example.com/"
					}
					return ""
				},
			}
			providedConfig := config.ProvideConfig("testing").(*bkapi.ClientConfig)

			Expect(providedConfig.Endpoint).To(Equal("http://dev-testing.example.com/"))
		})

		DescribeTable("should get app code from env", func(key string) {
			config := bkapi.ClientConfig{
				Getenv: func(k string) string {
					if k == key {
						return "app"
					}
					return ""
				},
			}

			providedConfig := config.ProvideConfig("testing").(*bkapi.ClientConfig)
			Expect(providedConfig.AppCode).To(Equal("app"))
		},
			Entry("BK_APP_CODE", "BK_APP_CODE"),
			Entry("APP_CODE", "APP_CODE"),
			Entry("BKPAAS_APP_ID", "BKPAAS_APP_ID"),
		)

		DescribeTable("should get app secret from env", func(key string) {
			config := bkapi.ClientConfig{
				Getenv: func(k string) string {
					if k == key {
						return "secret"
					}
					return ""
				},
			}

			providedConfig := config.ProvideConfig("testing").(*bkapi.ClientConfig)
			Expect(providedConfig.AppSecret).To(Equal("secret"))
		},
			Entry("BK_APP_SECRET", "BK_APP_SECRET"),
			Entry("SECRET_KEY", "SECRET_KEY"),
			Entry("BKPAAS_APP_SECRET", "BKPAAS_APP_SECRET"),
		)

		It("should return the client options", func() {
			option := mock.NewMockBkApiClientOption(ctrl)

			config := bkapi.ClientConfig{
				ClientOptions: []define.BkApiClientOption{option},
			}

			options := config.GetClientOptions()
			Expect(options[0]).To(Equal(option))
		})
	})
})
