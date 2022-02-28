package internal_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define/mock"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
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

			provider := internal.NewMarshalBodyProvider("text/plain", func(v interface{}) ([]byte, error) {
				return []byte(fmt.Sprintf("hello %s", v)), nil
			})
			Expect(provider.ProvideBody(operation, "world")).To(BeNil())
		})

		It("should marshal the data to json", func() {
			result := `{"hello":"world"}`

			operation.EXPECT().SetContentType("application/json").Return(operation)
			operation.EXPECT().SetContentLength(int64(len(result))).Return(operation)
			operation.EXPECT().SetBodyReader(gomock.Any()).DoAndReturn(func(reader io.Reader) define.Operation {
				body, err := ioutil.ReadAll(reader)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal(result))

				return operation
			})

			provider := internal.NewJsonBodyProvider()
			Expect(provider.ProvideBody(operation, map[string]string{
				"hello": "world",
			})).To(BeNil())
		})
	})

	Context("UnmarshalResultProvider", func() {
		It("should unmarshal the result", func() {
			var reader io.Reader

			provider := internal.NewUnmarshalResultProvider(func(body io.Reader, v interface{}) error {
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

		It("should marshal the data from json", func() {
			var result map[string]interface{}

			provider := internal.NewJsonResultProvider()
			Expect(provider.ProvideResult(&http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{"hello":"world"}`)),
			}, &result)).To(BeNil())

			Expect(result).To(Equal(map[string]interface{}{
				"hello": "world",
			}))
		})
	})
})
