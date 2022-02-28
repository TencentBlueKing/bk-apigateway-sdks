package internal_test

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define/mock"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provider", func() {
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

	Context("MarshalBodyProvider", func() {
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
})
