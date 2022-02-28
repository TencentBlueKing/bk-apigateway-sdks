package bkapi_test

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define/mock"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
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
			DoAndReturn(func(p define.BodyProvider) define.Operation {
				provider := p.(*internal.MarshalBodyProvider)
				Expect(provider.ContentType()).To(Equal("application/json"))
				return operation
			})

		option := bkapi.OptJsonBodyProvider()
		Expect(option.ApplyToOperation(operation)).To(Succeed())
	})
})
