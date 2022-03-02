package demo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/demo"
)

var _ = Describe("Client", func() {
	doSomethingWithOperation := func(op define.Operation) {
		// write your code here
	}

	It("should create a client by config", func() {
		client, err := demo.New(bkapi.ClientConfig{
			Endpoint:    "https://{api_name}.example.com/{stage}/",
			Stage:       "test",
			AccessToken: "access_token",
			AppCode:     "app_code",
			AppSecret:   "app_secret",
		})
		Expect(err).To(BeNil())

		Expect(client.Name()).To(Equal("demo"))

		doSomethingWithOperation(client.Anything())
	})
})
