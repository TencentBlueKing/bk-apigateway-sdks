package internal_test

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal/mock"
	"github.com/TencentBlueKing/gopkg/logging"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gentleman.v2"
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
			Expect(stacks).To(ContainElement(pluginA))
			Expect(stacks).To(ContainElement(pluginB))
		})

		It("should apply to client", func() {
			clientConfig := mock.NewMockClientConfig(ctrl)
			clientConfig.EXPECT().GetUrl().Return("").AnyTimes()
			clientConfig.EXPECT().GetAuthorizationHeaders().Return(nil).AnyTimes()
			clientConfig.EXPECT().GetLogger().Return(logging.GetLogger("")).AnyTimes()

			cli := internal.NewBkApiClient("", client, nil, clientConfig)

			Expect(cli.Apply(option)).To(Succeed())

			stacks := client.Middleware.GetStack()
			Expect(stacks).To(ContainElement(pluginA))
			Expect(stacks).To(ContainElement(pluginB))
		})
	})
})
