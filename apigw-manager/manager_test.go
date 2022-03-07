package manager_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	mgr "github.com/TencentBlueKing/bk-apigateway-sdks/apigw-manager"
	apigateway "github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-bk-apigateway"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gock.v1"
)

var _ = Describe("Manager", func() {
	var (
		manager *mgr.Manager
		apiName = "testing"
		config  bkapi.ClientConfig
	)

	BeforeEach(func() {
		config = bkapi.ClientConfig{
			Endpoint:  "http://example.com",
			AppCode:   "app_code",
			AppSecret: "app_secret",
		}

		var err error
		manager, err = mgr.NewManager(
			apiName, config, nil,
			func(configProvider define.ClientConfigProvider, opts ...define.BkApiClientOption) (*apigateway.Client, error) {
				opts = append(opts, bkapi.OptTransport(gock.NewTransport()))
				return apigateway.New(configProvider, opts...)
			},
		)
		Expect(err).To(BeNil())
	})

	It("should load definition from file", func() {
		dir, err := os.MkdirTemp("", "")
		Expect(err).To(BeNil())
		defer os.RemoveAll(dir)

		definitionFile := filepath.Join(dir, "test.yaml")
		Expect(ioutil.WriteFile(
			definitionFile,
			[]byte(`key: {{ data.value }}`),
			0644,
		)).To(Succeed())

		Expect(manager.LoadDefinition(definitionFile, map[string]interface{}{
			"value": "test",
		})).To(Succeed())

		definition := manager.GetDefinition()
		Expect(definition.Get("")).To(Equal(map[string]interface{}{
			"key": "test",
		}))
	})

	It("should return public key", func() {
		gock.New(config.Endpoint).
			Get("/api/v1/apis/testing/public_key/").
			Reply(200).
			JSON(map[string]interface{}{
				"code": 0,
				"data": map[string]interface{}{
					"public_key": "public_key",
				},
			})
		defer gock.Off()

		info, err := manager.GetPublicKey()
		Expect(err).To(BeNil())
		Expect(info).To(Equal(map[string]interface{}{
			"public_key": "public_key",
		}))
	})

	It("should return public key string", func() {
		gock.New(config.Endpoint).
			Get("/api/v1/apis/testing/public_key/").
			Reply(200).
			JSON(map[string]interface{}{
				"code": 0,
				"data": map[string]interface{}{
					"public_key": "public_key",
				},
			})
		defer gock.Off()

		key, err := manager.GetPublicKeyString()
		Expect(err).To(BeNil())
		Expect(key).To(Equal("public_key"))
	})
})
