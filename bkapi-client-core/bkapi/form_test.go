package bkapi_test

import (
	"io"
	"io/ioutil"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Form", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should provide urlencoded form", func() {
		result := "hello=world"
		operation := mock.NewMockOperation(ctrl)
		operation.EXPECT().SetContentType("application/x-www-form-urlencoded").Return(operation)
		operation.EXPECT().SetContentLength(int64(len(result))).Return(operation)
		operation.EXPECT().SetBodyReader(gomock.Any()).DoAndReturn(func(body io.Reader) define.Operation {
			data, err := ioutil.ReadAll(body)
			Expect(err).To(BeNil())
			Expect(string(data)).To(Equal(result))
			return operation
		})

		provider := bkapi.FormBodyProvider()
		Expect(provider.ProvideBody(operation, map[string][]string{
			"hello": {"world"},
		})).To(BeNil())
	})
})
