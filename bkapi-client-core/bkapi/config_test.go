package bkapi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
)

var _ = Describe("Config", func() {
	It("should clone a new config", func() {
		config := bkapi.Config{}
		Expect(config.GetName()).To(Equal(""))

		newConfig := config.Config("testing").(*bkapi.Config)
		Expect(newConfig.GetName()).To(Equal("testing"))

		Expect(config.GetName()).To(Equal(""))
	})

	It("should return endpoint as url", func() {
		config := bkapi.Config{
			Endpoint: "http://example.com",
		}

		Expect(config.Config("testing").GetUrl()).To(Equal("http://example.com"))
	})

	It("should render endpoint with params", func() {
		config := bkapi.Config{
			Endpoint: "http://{api_name}.example.com/{stage}/",
			Stage:    "prod",
		}

		Expect(config.Config("testing").GetUrl()).To(Equal("http://testing.example.com/prod/"))
	})

	It("should not return authorization headers when related params are empty", func() {
		config := bkapi.Config{
			AccessToken:         "",
			AuthorizationJWT:    "",
			AppCode:             "",
			AppSecret:           "",
			AuthorizationParams: nil,
		}

		Expect(config.GetAuthorizationHeaders()).To(BeEmpty())
	})

	It("should return access token authorization headers", func() {
		config := bkapi.Config{
			AccessToken:      "access_token",
			AuthorizationJWT: "jwt",
		}

		Expect(config.GetAuthorizationHeaders()).To(Equal(map[string]string{
			"X-Bkapi-Authorization": `{"access_token":"access_token","jwt":"jwt"}`,
		}))
	})

	It("should return app code authorization headers", func() {
		config := bkapi.Config{
			AppCode:   "app_code",
			AppSecret: "app_secret",
		}

		Expect(config.GetAuthorizationHeaders()).To(Equal(map[string]string{
			"X-Bkapi-Authorization": `{"bk_app_code":"app_code","bk_app_secret":"app_secret"}`,
		}))
	})

	It("should return common authorization headers", func() {
		config := bkapi.Config{
			AuthorizationParams: map[string]string{
				"bk_token": "token",
			},
		}

		Expect(config.GetAuthorizationHeaders()).To(Equal(map[string]string{
			"X-Bkapi-Authorization": `{"bk_token":"token"}`,
		}))
	})

	It("should return authorization headers marshal by custom marshaler", func() {
		config := bkapi.Config{
			AccessToken: "access_token",
			JsonMarshaler: func(v interface{}) ([]byte, error) {
				return []byte(`{"access_token": "access_token"}`), nil
			},
		}

		Expect(config.GetAuthorizationHeaders()).To(Equal(map[string]string{
			"X-Bkapi-Authorization": `{"access_token": "access_token"}`,
		}))
	})
})
