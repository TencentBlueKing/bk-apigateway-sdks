package bkapi_test

import (
	"net/http"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/transport"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define/mock"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
)

var _ = Describe("Multipart", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mock.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		roundTripper = mock.NewMockRoundTripper(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should set json body", func() {
		client := gentleman.New()
		client.Use(transport.Set(roundTripper))

		request := client.Request()
		operation := internal.NewOperation("", request)

		option := bkapi.OptMultipartFormBodyProvider()
		Expect(option.ApplyToOperation(operation)).To(Succeed())

		roundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			Expect(req.Header.Get("Content-Type")).To(HavePrefix("multipart/form-data"))

			return &http.Response{}, nil
		})

		_, err := operation.
			SetBody(map[string][]string{
				"hello": {"world"},
			}).
			Request()
		Expect(err).To(BeNil())
	})
})
