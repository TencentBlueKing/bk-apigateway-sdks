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
		ctrl                    *gomock.Controller
		client                  *internal.BkApiClient
		gentlemanClient         *gentleman.Client
		operation               *mock.MockOperation
		request                 *gentleman.Request
		mockTransport           *mock.MockRoundTripper
		operationConfig         *mock.MockOperationConfig
		operationConfigProvider *mock.MockOperationConfigProvider
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
				operation.EXPECT().Name().Return(name).AnyTimes()

				request = req
				return operation
			},
		)
		operationConfig = mock.NewMockOperationConfig(ctrl)
		operationConfig.EXPECT().GetName().Return("").AnyTimes()
		operationConfig.EXPECT().GetPath().Return("").AnyTimes()
		operationConfig.EXPECT().GetMethod().Return("").AnyTimes()

		operationConfigProvider = mock.NewMockOperationConfigProvider(ctrl)
		operationConfigProvider.EXPECT().ProvideConfig().Return(operationConfig).AnyTimes()
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
			operationConfig := mock.NewMockOperationConfig(ctrl)
			operationConfigProvider := mock.NewMockOperationConfigProvider(ctrl)
			operationConfigProvider.EXPECT().ProvideConfig().Return(operationConfig).AnyTimes()

			operationConfig.EXPECT().GetName().Return("test").AnyTimes()
			operationConfig.EXPECT().GetPath().Return("/test").AnyTimes()
			operationConfig.EXPECT().GetMethod().Return("POST").AnyTimes()

			operation.EXPECT().Apply(gomock.Any()).AnyTimes()

			op := client.NewOperation(operationConfigProvider)
			Expect(op).To(Equal(operation))

			mockTransportRoundTrip()
			// make a request to apply the middlewares
			_, err := request.Send()
			Expect(err).To(BeNil())

			Expect(request.Context.Request.Method).To(Equal("POST"))
			Expect(request.Context.Request.URL.Path).To(Equal("/test"))
		})

		It("should new an operation with option", func() {
			option := mock.NewMockOperationOption(ctrl)
			operation.EXPECT().Apply(option)

			op := client.NewOperation(operationConfigProvider, option)
			Expect(op).NotTo(BeNil())
		})

		It("should new an operation with common operation option", func() {
			option := mock.NewMockOperationOption(ctrl)
			operation.EXPECT().Apply(option)
			Expect(client.AddOperationOptions(option)).To(BeNil())

			op := client.NewOperation(operationConfigProvider)
			Expect(op).NotTo(BeNil())
		})

		It("should generate operation name by config.Name", func() {
			operationConfig := mock.NewMockOperationConfig(ctrl)
			operationConfigProvider := mock.NewMockOperationConfigProvider(ctrl)
			operationConfigProvider.EXPECT().ProvideConfig().Return(operationConfig).AnyTimes()

			operationConfig.EXPECT().GetName().Return("operation").AnyTimes()
			operationConfig.EXPECT().GetMethod().Return("GET").AnyTimes()
			operationConfig.EXPECT().GetPath().Return("/test").AnyTimes()

			operation := client.NewOperation(operationConfigProvider)

			Expect(operation.Name()).To(Equal("testing.operation"))
		})

		It("should generate anonymous operation name", func() {
			operationConfig := mock.NewMockOperationConfig(ctrl)
			operationConfigProvider := mock.NewMockOperationConfigProvider(ctrl)
			operationConfigProvider.EXPECT().ProvideConfig().Return(operationConfig).AnyTimes()

			operationConfig.EXPECT().GetName().Return("").AnyTimes()
			operationConfig.EXPECT().GetMethod().Return("GET").AnyTimes()
			operationConfig.EXPECT().GetPath().Return("/test").AnyTimes()

			operation := client.NewOperation(operationConfigProvider)

			Expect(operation.Name()).To(Equal("testing(GET /test)"))
		})

		It("should set the user agent", func() {
			op := client.NewOperation(operationConfigProvider)
			Expect(op).To(Equal(operation))

			mockTransportRoundTrip()

			_, err := request.Send()
			Expect(err).To(BeNil())

			request := request.Context.Request
			Expect(request.Header.Get("User-Agent")).To(Equal(fmt.Sprintf("%s/%s", define.UserAgent, define.Version)))
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
