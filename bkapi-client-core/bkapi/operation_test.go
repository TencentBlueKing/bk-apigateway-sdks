package bkapi_test

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Operation", func() {
	Context("OperationConfig", func() {
		It("should clone a config", func() {
			config := bkapi.OperationConfig{
				Name:   "test",
				Path:   "/test",
				Method: "GET",
			}

			providedConfig := config.ProvideConfig()

			config.Name = ""
			config.Path = ""
			config.Method = ""

			Expect(providedConfig.GetName()).To(Equal("test"))
			Expect(providedConfig.GetPath()).To(Equal("/test"))
			Expect(providedConfig.GetMethod()).To(Equal("GET"))
		})
	})

	Context("OperationOption", func() {
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

	Context("Option", func() {
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

		It("should set the body", func() {
			body := make(map[string]interface{})
			operation.EXPECT().SetBody(body).Return(nil)

			option := bkapi.OptSetOperationBody(body)
			Expect(option.ApplyToOperation(operation)).To(Succeed())
		})

		It("should set the result", func() {
			result := make(map[string]interface{})
			operation.EXPECT().SetResult(result).Return(nil)

			option := bkapi.OptSetOperationResult(result)
			Expect(option.ApplyToOperation(operation)).To(Succeed())
		})
	})
})
