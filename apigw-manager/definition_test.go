package manager_test

import (
	manager "github.com/TencentBlueKing/bk-apigateway-sdks/apigw-manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Definition", func() {
	Context("Get", func() {
		var definition *manager.Definition

		BeforeEach(func() {
			definition = manager.NewDefinition(map[string]interface{}{
				"sub": map[string]interface{}{
					"name": "testing",
					"value": map[string]interface{}{
						"foo": "bar",
					},
				},
			})
		})

		It("should return a subdefinition", func() {
			sub, err := definition.Get("sub")
			Expect(err).To(BeNil())
			Expect(sub["name"]).To(Equal("testing"))
		})

		It("should return a deep subdefinition", func() {
			sub, err := definition.Get("sub.value")
			Expect(err).To(BeNil())
			Expect(sub["foo"]).To(Equal("bar"))
		})

		It("should return an error when subdefinition not found", func() {
			_, err := definition.Get("not_found")
			Expect(err).ToNot(BeNil())
		})
	})

	It("should new definition from yaml", func() {
		definition, err := manager.NewDefinitionFromYaml([]byte("sub:\n  name: testing"))
		Expect(err).To(BeNil())
		sub, err := definition.Get("sub")
		Expect(err).To(BeNil())
		Expect(sub["name"]).To(Equal("testing"))
	})
})
