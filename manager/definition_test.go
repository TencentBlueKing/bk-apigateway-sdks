/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2017 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package manager_test

import (
	manager "github.com/TencentBlueKing/bk-apigateway-sdks/manager"
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
