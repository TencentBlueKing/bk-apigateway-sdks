package demo_test

import (
	"net/http"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/demo"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Demo", func() {
	var (
		ctrl          *gomock.Controller
		mockServer    *mock.MockRoundTripper
		mockServerOpt define.BkapiOption
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockServer = mock.NewMockRoundTripper(ctrl)
		mockServerOpt = bkapi.OptTransport(mockServer)

		mockServer.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Request:       req,
				Header:        req.Header,
				ContentLength: req.ContentLength,
				StatusCode:    200,
				Body:          req.Body,
			}, nil
		}).AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should create a client by config", func() {
		client, err := demo.New(demo.Config{
			Endpoint:    "https://{api_name}.example.com/{stage}/",
			Stage:       "test",
			AccessToken: "access_token",
		}, mockServerOpt)
		Expect(err).To(BeNil())

		Expect(client.Name()).To(Equal("demo"))

		response, err := client.Echo().Request()
		Expect(err).To(BeNil())

		request := response.Request
		Expect(request.URL.String()).To(Equal("https://demo.example.com/test/echo/"))
		Expect(request.Header.Get("X-Bkapi-Authorization")).To(Equal(`{"access_token":"access_token"}`))
	})

	It("should request api with chaining call", func() {
		client, err := demo.New(demo.Config{
			Endpoint: "https://example.com/",
		}, mockServerOpt)

		Expect(err).To(BeNil())

		result := make(map[string]interface{})
		response, err := client.Echo().
			SetBodyProvider(bkapi.JsonBodyProvider()).
			SetBody(map[string]interface{}{
				"from": "body",
			}).
			SetResultProvider(bkapi.JsonResultProvider()).
			SetResult(&result).
			SetQueryParams(map[string]string{
				"from": "query",
			}).
			SetHeaders(map[string]string{
				"X-Header": "my-header",
			}).
			Request()
		Expect(err).To(BeNil())

		Expect(result["from"]).To(Equal("body"))
		Expect(response.Header.Get("X-Header")).To(Equal("my-header"))
		Expect(response.Request.URL.Query().Get("from")).To(Equal("query"))
	})

	It("should request api with optional call", func() {
		client, err := demo.New(demo.Config{
			Endpoint: "https://example.com/",
		}, mockServerOpt)

		Expect(err).To(BeNil())

		result := make(map[string]interface{})
		response, err := client.Echo(
			bkapi.OptJsonBodyProvider(),
			bkapi.OptSetOperationBody(map[string]interface{}{
				"from": "body",
			}),
			bkapi.OptJsonResultProvider(),
			bkapi.OptSetOperationResult(&result),
			bkapi.OptSetQueryParams(map[string]string{
				"from": "query",
			}),
			bkapi.OptSetHeaders(map[string]string{
				"X-Header": "my-header",
			}),
		).Request()
		Expect(err).To(BeNil())

		Expect(result["from"]).To(Equal("body"))
		Expect(response.Header.Get("X-Header")).To(Equal("my-header"))
		Expect(response.Request.URL.Query().Get("from")).To(Equal("query"))
	})

	It("should apply common operation options", func() {
		client, err := demo.New(demo.Config{
			Endpoint: "https://example.com/",
		},
			mockServerOpt,
			bkapi.OptJsonBodyProvider(),
			bkapi.OptJsonResultProvider(),
			bkapi.OptSetHeaders(map[string]string{
				"X-Header": "my-header",
			}),
		)

		Expect(err).To(BeNil())

		result := make(map[string]interface{})
		response, err := client.Echo().
			SetBody(map[string]interface{}{
				"from": "body",
			}).
			SetResult(&result).
			Request()
		Expect(err).To(BeNil())

		Expect(result["from"]).To(Equal("body"))
		Expect(response.Header.Get("X-Header")).To(Equal("my-header"))
	})

	It("should call a api dynamically", func() {
		client, err := demo.New(demo.Config{
			Endpoint: "https://example.com/",
		}, mockServerOpt)

		Expect(err).To(BeNil())

		response, err := client.NewOperation(define.OperationConfig{
			Method: "GET",
			Path:   "/thing/{id}",
		}).
			SetPathParams(map[string]string{
				"id": "1",
			}).
			Request()

		Expect(err).To(BeNil())
		Expect(response.Request.URL.String()).To(Equal("https://example.com/thing/1"))
	})
})
