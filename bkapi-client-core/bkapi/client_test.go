package bkapi_test

import (
	"net/http"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Context("BkApiClient", func() {
		var (
			ctrl                 *gomock.Controller
			apiName              = "testing"
			configProvider       *mock.MockClientConfigProvider
			config               *mock.MockClientConfig
			roundTripper         *mock.MockRoundTripper
			roundTripperOpt      define.BkapiOption
			url                  string
			authorizationHeaders map[string]string
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
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

		It("should apply option", func() {
			client, err := bkapi.NewBkApiClient(apiName, configProvider, roundTripperOpt)
			Expect(err).To(BeNil())

			request := getMockRequest(client)
			Expect(request).NotTo(BeNil())
		})

		It("should apply url from config", func() {
			url = "http://example.com/"
			client, err := bkapi.NewBkApiClient(apiName, configProvider, roundTripperOpt)
			Expect(err).To(BeNil())

			request := getMockRequest(client)
			Expect(request.URL.String()).To(Equal(url))
		})

		It("should set the authorization header", func() {
			authorizationHeaders["Authorization"] = "Bearer token"

			client, err := bkapi.NewBkApiClient(apiName, configProvider, roundTripperOpt)
			Expect(err).To(BeNil())

			request := getMockRequest(client)
			Expect(request.Header.Get("Authorization")).To(Equal("Bearer token"))
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
				Endpoint: "http://{api_name}.example.com/{stage}/",
				Stage:    "prod",
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
			}))
		})

		It("should return app code authorization headers", func() {
			config := bkapi.ClientConfig{
				AppCode:   "app_code",
				AppSecret: "app_secret",
			}

			Expect(config.GetAuthorizationHeaders()).To(Equal(map[string]string{
				"X-Bkapi-Authorization": `{"bk_app_code":"app_code","bk_app_secret":"app_secret"}`,
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
			}))
		})

		It("should set stage by default", func() {
			config := bkapi.ClientConfig{}
			providedConfig := config.ProvideConfig("testing").(*bkapi.ClientConfig)

			Expect(providedConfig.Stage).To(Equal("prod"))
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

		It("should set endpoint by env BK_API_URL_TMPL", func() {
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
		)
	})
})
