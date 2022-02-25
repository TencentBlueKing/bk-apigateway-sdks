package internal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
)

var _ = Describe("Placeholder", func() {
	Context("ReplacePlaceHolder", func() {
		It("should replace the placeholder with the given string", func() {
			Expect(internal.ReplacePlaceHolder(
				"{param}", map[string]string{"param": "value"},
			)).To(Equal("value"))
		})

		It("should replace the placeholder with spaces", func() {
			Expect(internal.ReplacePlaceHolder(
				"{ param }", map[string]string{"param": "value"},
			)).To(Equal("value"))
		})

		It("should not replace the placeholder when param is not found", func() {
			Expect(internal.ReplacePlaceHolder(
				"{param}", map[string]string{},
			)).To(Equal("{param}"))
		})

		It("should replace the placeholders", func() {
			Expect(internal.ReplacePlaceHolder(
				"{action} {name}", map[string]string{
					"action": "hello",
					"name":   "world",
				},
			)).To(Equal("hello world"))
		})
	})
})
