package internal_test

import (
	"fmt"
	"net/http"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/transport"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal/mock"
)

var _ = Describe("Client", func() {
	var (
		ctrl            *gomock.Controller
		client          *internal.BkApiClient
		gentlemanClient *gentleman.Client
		operation       *mock.MockOperation
		request         *gentleman.Request
		mockTransport   *mock.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockTransport = mock.NewMockRoundTripper(ctrl)
		operation = mock.NewMockOperation(ctrl)

		gentlemanClient = gentleman.New()
		gentlemanClient.Use(transport.Set(mockTransport))

		client = internal.NewBkApiClient(
			"testing", gentlemanClient,
			func(name string, req *gentleman.Request) define.Operation {
				request = req
				return operation
			},
		)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	mockTransportRoundTrip := func() {
		mockTransport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Request: req,
			}, nil
		})
	}

	Context("BkApiClient", func() {
		AfterEach(func() {
			ctrl.Finish()
		})

		It("should fail on apply", func() {
			option := mock.NewMockBkApiClientOption(ctrl)

			option.EXPECT().ApplyToClient(gomock.Any()).Return(errors.New("test"))
			Expect(client.Apply(option)).NotTo(BeNil())
		})

		It("should apply an option", func() {
			option := mock.NewMockBkApiClientOption(ctrl)
			option.EXPECT().ApplyToClient(client).Return(nil)

			Expect(client.Apply(option)).To(BeNil())
		})

		It("should new an operation", func() {
			operation.EXPECT().Apply(gomock.Any()).AnyTimes()

			op := client.NewOperation(define.OperationConfig{
				Method: "POST",
				Path:   "/test",
			})
			Expect(op).To(Equal(operation))

			mockTransportRoundTrip()
			// make a request to apply the middlewares
			_, err := request.Send()
			Expect(err).To(BeNil())

			Expect(request.Context.Request.Method).To(Equal("POST"))
			Expect(request.Context.Request.URL.Path).To(Equal("test"))
		})

		It("should new an operation with options", func() {
			option := mock.NewMockOperationOption(ctrl)
			operation.EXPECT().Apply(option)

			op := client.NewOperation(define.OperationConfig{}, option)
			Expect(op).NotTo(BeNil())
		})
	})

	Context("BkApiClientOption", func() {
		It("should fail when the type is not supported", func() {
			var client mock.MockBkApiClient

			opt := internal.NewBkApiClientOption(nil)
			err := opt.ApplyToClient(&client)

			Expect(errors.Cause(err)).To(Equal(define.ErrTypeNotMatch))
		})

		It("should apply function", func() {
			var client internal.BkApiClient
			err := fmt.Errorf("testing")

			opt := internal.NewBkApiClientOption(func(c *internal.BkApiClient) error {
				Expect(c).To(Equal(&client))
				return err
			})

			Expect(opt.ApplyToClient(&client)).To(Equal(err))
		})
	})
})
