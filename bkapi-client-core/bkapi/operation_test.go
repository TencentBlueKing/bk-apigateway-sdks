package bkapi_test

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Operation", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should apply to client", func() {
		option := bkapi.NewOperationOption(nil)

		client := mock.NewMockBkApiClient(ctrl)
		client.EXPECT().AddOperationOptions(option).Return(nil)

		Expect(option.ApplyToClient(client)).To(Succeed())
	})

	It("should apply to operation", func() {
		operation := mock.NewMockOperation(ctrl)

		called := false
		option := bkapi.NewOperationOption(func(op define.Operation) error {
			called = true
			Expect(op).To(Equal(operation))

			return nil
		})

		Expect(option.ApplyToOperation(operation)).To(Succeed())
		Expect(called).To(BeTrue())
	})
})
