package bkapi_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provider", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("MarshalBodyProvider", func() {
		var operation *mock.MockOperation

		BeforeEach(func() {
			operation = mock.NewMockOperation(ctrl)
		})

		It("should marshal the data", func() {
			result := "hello world"

			operation.EXPECT().SetContentType("text/plain").Return(operation)
			operation.EXPECT().SetContentLength(int64(len(result))).Return(operation)
			operation.EXPECT().SetBodyReader(gomock.Any()).DoAndReturn(func(reader io.Reader) define.Operation {
				body, err := ioutil.ReadAll(reader)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal(result))

				return operation
			})

			provider := bkapi.NewMarshalBodyProvider("text/plain", func(v interface{}) ([]byte, error) {
				return []byte(fmt.Sprintf("hello %s", v)), nil
			})
			Expect(provider.ProvideBody(operation, "world")).To(BeNil())
		})
	})

	Context("UnmarshalResultProvider", func() {
		It("should unmarshal the result", func() {
			var reader io.Reader

			provider := bkapi.NewUnmarshalResultProvider(func(body io.Reader, v interface{}) error {
				*v.(*io.Reader) = body

				return nil
			})

			Expect(provider.ProvideResult(&http.Response{
				Body: ioutil.NopCloser(strings.NewReader("hello world")),
			}, &reader)).To(BeNil())

			result, err := ioutil.ReadAll(reader)
			Expect(err).To(BeNil())
			Expect(string(result)).To(Equal("hello world"))
		})
	})

	DescribeTable("provider as BkApiClientOption", func(provider interface{}) {
		opt, ok := provider.(define.BkApiClientOption)
		Expect(ok).To(BeTrue())

		client := mock.NewMockBkApiClient(ctrl)
		client.EXPECT().AddOperationOptions(opt).Return(nil)

		Expect(opt.ApplyToClient(client)).To(BeNil())
	},
		Entry("MarshalBodyProvider", bkapi.NewMarshalBodyProvider("", nil)),
		Entry("UnmarshalResultProvider", bkapi.NewUnmarshalResultProvider(nil)),
		Entry("FunctionalBodyProvider", bkapi.NewFunctionalBodyProvider(nil)),
	)

	DescribeTable("BodyProvider as OperationOption", func(provider define.BodyProvider) {
		opt, ok := provider.(define.OperationOption)
		Expect(ok).To(BeTrue())

		operation := mock.NewMockOperation(ctrl)
		operation.EXPECT().SetBodyProvider(provider).Return(nil)

		Expect(opt.ApplyToOperation(operation)).To(BeNil())
	},
		Entry("MarshalBodyProvider", bkapi.NewMarshalBodyProvider("", nil)),
		Entry("FunctionalBodyProvider", bkapi.NewFunctionalBodyProvider(nil)),
	)

	DescribeTable("ResultProvider as OperationOption", func(provider define.ResultProvider) {
		opt, ok := provider.(define.OperationOption)
		Expect(ok).To(BeTrue())

		operation := mock.NewMockOperation(ctrl)
		operation.EXPECT().SetResultProvider(provider).Return(nil)

		Expect(opt.ApplyToOperation(operation)).To(BeNil())
	},
		Entry("UnmarshalResultProvider", bkapi.NewUnmarshalResultProvider(nil)),
	)
})