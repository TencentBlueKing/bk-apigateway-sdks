package demo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/demo"
)

var _ = Describe("Operation", func() {
	var (
		client *demo.Client
	)

	BeforeEach(func() {
		var err error
		client, err = demo.New(bkapi.ClientConfig{
			Endpoint: "http://httpbin.org/",
		})

		Expect(err).To(BeNil())
	})

	Context("Example for common usage", func() {
		It("should request to anything", func() {
			response, err := client.Anything().
				Request()

			// you can handle the response here
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
		})

		Context("decode result", func() {
			It("request to anything by chaining style", func() {
				var result demo.AnythingResponse

				_, _ = client.Anything().
					SetResultProvider(bkapi.JsonResultProvider()).
					SetResult(&result).
					Request()

				// when you has set the result provider, the result will be decoded automatically
				Expect(result.URL).To(Equal("http://httpbin.org/anything"))
			})

			It("request to anything by option style", func() {
				var result demo.AnythingResponse

				_, _ = client.Anything(
					bkapi.OptJsonResultProvider(),
					bkapi.OptSetRequestResult(&result),
				).Request()

				// it is also ok to set the result provider by option style
				Expect(result.URL).To(Equal("http://httpbin.org/anything"))
			})
		})

		Context("with request payload", func() {
			var (
				result *demo.AnythingResponse
			)

			BeforeEach(func() {
				result = new(demo.AnythingResponse)

				// you can set the common options to the client
				err := client.AddOperationOptions(
					bkapi.JsonBodyProvider(),
					bkapi.JsonResultProvider(),
					bkapi.OptSetRequestResult(result),
				)

				Expect(err).To(BeNil())
			})

			It("request to anything by chaining style", func() {
				_, _ = client.Anything().
					SetQueryParams(map[string]string{
						"from": "query",
					}).
					SetBody(map[string]interface{}{
						"from": "body",
					}).
					SetHeaders(map[string]string{
						"X-Header": "my-header",
					}).
					Request()

				Expect(result.Args["from"]).To(Equal("query"))
				Expect(result.JSON["from"]).To(Equal("body"))
				Expect(result.Headers["X-Header"]).To(Equal("my-header"))
			})

			It("request to anything by option style", func() {
				_, _ = client.Anything(
					bkapi.OptSetRequestQueryParams(map[string]string{
						"from": "query",
					}),
					bkapi.OptSetRequestBody(map[string]interface{}{
						"from": "body",
					}),
					bkapi.OptSetRequestHeaders(map[string]string{
						"X-Header": "my-header",
					}),
				).Request()

				Expect(result.Args["from"]).To(Equal("query"))
				Expect(result.JSON["from"]).To(Equal("body"))
				Expect(result.Headers["X-Header"]).To(Equal("my-header"))
			})
		})

		Context("with path params", func() {
			It("request to anything by chaining style", func() {
				response, _ := client.StatusCode().
					SetPathParams(map[string]string{
						"code": "200",
					}).
					Request()

				Expect(response.Request.URL.Path).To(Equal("/status/200"))
			})

			It("request to anything by option style", func() {
				response, _ := client.StatusCode(
					bkapi.OptSetRequestPathParams(map[string]string{
						"code": "200",
					}),
				).Request()

				Expect(response.Request.URL.Path).To(Equal("/status/200"))
			})
		})
	})

	Context("Example for error handling", func() {
		It("should handle 5XX", func() {
			response, err := client.StatusCode().
				SetPathParams(map[string]string{
					"code": "500",
				}).
				Request()

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(500))
		})

		It("should handle decode error", func() {
			var result map[string]interface{}

			// 5xx response will sent non-json body
			_, err := client.StatusCode().
				SetResultProvider(bkapi.JsonResultProvider()).
				SetResult(&result).
				SetPathParams(map[string]string{
					"code": "500",
				}).
				Request()

			Expect(err).NotTo(BeNil())
		})
	})
})
