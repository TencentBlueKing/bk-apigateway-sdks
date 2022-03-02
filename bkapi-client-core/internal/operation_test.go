package internal_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/transport"
)

var _ = Describe("operation", func() {
	Context("Operation", func() {
		var (
			ctrl          *gomock.Controller
			mockTransport *mock.MockRoundTripper
			client        *gentleman.Client
			response      *http.Response
			operation     *internal.Operation
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockTransport = mock.NewMockRoundTripper(ctrl)

			client = gentleman.New()
			client.Use(transport.Set(mockTransport))

			response = &http.Response{}

			request := client.Request()
			operation = internal.NewOperation("test", request)
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		mockTransportRoundTrip := func() {
			mockTransport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
				response.Request = req
				return response, nil
			})
		}

		It("should fail on apply", func() {
			option := mock.NewMockOperationOption(ctrl)
			option.EXPECT().ApplyToOperation(gomock.Any()).Return(errors.New("test"))

			_, err := operation.
				Apply(option).
				Request()

			Expect(err).NotTo(BeNil())
		})

		It("should fail on roundtrip", func() {
			mockTransport.EXPECT().RoundTrip(gomock.Any()).Return(nil, fmt.Errorf("testing"))

			response, err := operation.Request()
			Expect(err).NotTo(BeNil())
			Expect(response).To(BeNil())
		})

		It("should set request headers", func() {
			mockTransportRoundTrip()

			response, err := operation.
				SetHeaders(map[string]string{
					"Content-Type": "application/json",
				}).
				Request()
			Expect(err).To(BeNil())

			Expect(response.Request.Header.Get("Content-Type")).To(Equal("application/json"))
		})

		It("should set request params", func() {
			mockTransportRoundTrip()

			response, err := operation.
				SetQueryParams(map[string]string{
					"foo": "bar",
				}).
				Request()
			Expect(err).To(BeNil())

			Expect(response.Request.URL.Query().Get("foo")).To(Equal("bar"))
		})

		It("should set request path params", func() {
			mockTransportRoundTrip()

			request := client.Request().Path("/hello/{name}")
			operation = internal.NewOperation("test", request)

			response, err := operation.
				SetPathParams(map[string]string{
					"name": "world",
				}).
				Request()
			Expect(err).To(BeNil())

			Expect(response.Request.URL.Path).To(Equal("/hello/world"))
		})

		It("should set request body", func() {
			mockTransportRoundTrip()

			response, err := operation.
				SetBodyReader(strings.NewReader("testing")).
				Request()
			Expect(err).To(BeNil())

			body, err := ioutil.ReadAll(response.Request.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("testing"))
		})

		It("should set request body with json", func() {
			requestBody := []byte(`{"foo":"bar"}`)

			provider := mock.NewMockBodyProvider(ctrl)
			provider.EXPECT().ProvideBody(operation, gomock.Any()).DoAndReturn(func(op define.Operation, data interface{}) error {
				op.
					SetContentType("application/json").
					SetContentLength(int64(len(requestBody))).
					SetBodyReader(bytes.NewReader(requestBody))

				return nil
			})

			mockTransportRoundTrip()

			response, err := operation.
				SetBodyProvider(provider).
				SetBody(gomock.Any()).
				Request()

			Expect(err).To(BeNil())
			Expect(response.Request.Header.Get("Content-Type")).To(Equal("application/json"))
			Expect(int(response.Request.ContentLength)).To(Equal(len(requestBody)))

			body, err := ioutil.ReadAll(response.Request.Body)
			Expect(err).To(BeNil())
			Expect(body).To(Equal(requestBody))
		})

		It("should decode response body", func() {
			mockTransportRoundTrip()

			provider := mock.NewMockResultProvider(ctrl)
			provider.EXPECT().ProvideResult(gomock.Any(), gomock.Any()).DoAndReturn(func(response *http.Response, r interface{}) error {
				result := r.(map[string]interface{})
				result["foo"] = "bar"

				return nil
			})
			result := make(map[string]interface{})

			_, err := operation.
				SetResult(result).
				SetResultProvider(provider).
				Request()

			Expect(err).To(BeNil())
			Expect(result["foo"]).To(Equal("bar"))
		})

		It("should set request context", func() {
			ctx := context.WithValue(context.Background(), "key", "testing")

			mockTransportRoundTrip()

			response, err := operation.SetContext(ctx).Request()
			Expect(err).To(BeNil())

			requestCtx := response.Request.Context()
			Expect(requestCtx.Value("key")).To(Equal("testing"))
		})

		It("should close response body", func() {
			body := mock.NewMockReadCloser(ctrl)
			body.EXPECT().Read(gomock.Any()).Return(0, io.EOF).AnyTimes()
			body.EXPECT().Close().Return(nil)
			response.Body = body

			mockTransportRoundTrip()

			_, err := operation.Request()
			Expect(err).To(BeNil())
			Expect(response.Close).To(BeTrue())
		})
	})

	Context("OperationOption", func() {
		It("should fail when the operation type is not supported", func() {
			var mockOperation mock.MockOperation

			opt := internal.NewOperationOption(nil)
			err := opt.ApplyToOperation(&mockOperation)

			Expect(errors.Cause(err)).To(Equal(define.ErrTypeNotMatch))
		})

		It("should apply function", func() {
			var operation internal.Operation
			err := fmt.Errorf("testing")

			opt := internal.NewOperationOption(func(o *internal.Operation) error {
				Expect(o).To(Equal(&operation))
				return err
			})

			Expect(opt.ApplyToOperation(&operation)).To(Equal(err))
		})
	})
})
