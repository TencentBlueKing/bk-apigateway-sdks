package internal_test

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	dmock "github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define/mock"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

var _ = Describe("Option", func() {
	Context("PluginOption", func() {
		var (
			ctrl             *gomock.Controller
			option           *internal.PluginOption
			pluginA, pluginB *mock.MockPlugin
			client           *gentleman.Client
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			pluginA = mock.NewMockPlugin(ctrl)
			pluginB = mock.NewMockPlugin(ctrl)
			option = internal.NewPluginOption(pluginA, pluginB)
			client = gentleman.New()
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should apply to operation", func() {
			request := client.Request()
			operation := internal.NewOperation("", request)
			operation.Apply(option)
			Expect(operation.GetError()).To(BeNil())

			stacks := request.Middleware.GetStack()
			Expect(stacks).To(Equal([]plugin.Plugin{pluginA, pluginB}))
		})

		It("should apply to client", func() {
			operation := dmock.NewMockOperation(ctrl)
			operation.EXPECT().Apply(option).Return(nil)

			client := internal.NewBkApiClient("", client, func(name string, request *gentleman.Request) define.Operation {
				return operation
			})
			Expect(client.Apply(option)).To(Succeed())

			client.NewOperation(define.OperationConfig{})
		})
	})
})