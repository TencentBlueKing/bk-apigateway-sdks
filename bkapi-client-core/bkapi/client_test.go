package bkapi_test

import (
	"net/http"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
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

		configProvider.EXPECT().Config(apiName).Return(config).AnyTimes()
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

		operation := client.NewOperation(define.OperationConfig{})
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
