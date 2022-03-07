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
					apiName, config, nil,
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
