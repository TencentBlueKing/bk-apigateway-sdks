package bkapi_test

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Multipart", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should provide multopart form", func() {
		operation := mock.NewMockOperation(ctrl)
		operation.EXPECT().Apply(gomock.Any()).Return(operation)

		provider := bkapi.MultipartFormBodyProvider()
		Expect(provider.ProvideBody(operation, map[string][]string{
			"hello": {"world"},
		})).To(BeNil())
	})
})
