package bkapi_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Json", func() {
	var (
		ctrl      *gomock.Controller
		operation *mock.MockOperation
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		operation = mock.NewMockOperation(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should set json body", func() {
		operation.EXPECT().SetBodyProvider(gomock.Any()).
			DoAndReturn(func(provider define.BodyProvider) define.Operation {
				result := `{"hello":"world"}`

				operation.EXPECT().SetContentType("application/json").Return(operation)
				operation.EXPECT().SetContentLength(int64(len(result))).Return(operation)
				operation.EXPECT().SetBodyReader(gomock.Any()).DoAndReturn(func(reader io.Reader) define.Operation {
					body, err := ioutil.ReadAll(reader)
					Expect(err).To(BeNil())
					Expect(string(body)).To(Equal(result))

					return operation
				})

				Expect(provider.ProvideBody(operation, map[string]string{
					"hello": "world",
				})).To(BeNil())

				return operation
			})

		option := bkapi.OptJsonBodyProvider()
		Expect(option.ApplyToOperation(operation)).To(Succeed())
	})

	It("should set json result", func() {
		operation.EXPECT().SetResultProvider(gomock.Any()).
			DoAndReturn(func(provider define.ResultProvider) define.Operation {
				var result map[string]interface{}

				Expect(provider.ProvideResult(&http.Response{
					Body: ioutil.NopCloser(strings.NewReader(`{"hello":"world"}`)),
				}, &result)).To(BeNil())

				Expect(result).To(Equal(map[string]interface{}{
					"hello": "world",
				}))

				return operation
			})

		option := bkapi.OptJsonResultProvider()
		Expect(option.ApplyToOperation(operation)).To(Succeed())
	})
})
