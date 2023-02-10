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

package demo_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/demo"
)

var _ = Describe("Operation", func() {
	var (
		client *demo.Client
	)

	BeforeEach(func() {
		var err error
		client, err = demo.New(bkapi.ClientConfig{
			Endpoint: "https://httpbin.org/",
		})

		Expect(err).To(BeNil())
	})

	Context("Example for common usage", func() {
		It("should request to anything", func() {
			response, err := client.Anything().
				Request()

			if err != nil {
				Skip(err.Error())
			}

			// you can handle the response here
			Expect(response.StatusCode).To(Equal(200))
		})

		Context("decode result", func() {
			It("request to anything by chaining style", func() {
				var result demo.AnythingResponse

				_, err := client.Anything().
					SetResultProvider(bkapi.JsonResultProvider()).
					SetResult(&result).
					Request()

				if err != nil {
					Skip(err.Error())
				}

				// when you has set the result provider, the result will be decoded automatically
				Expect(result.URL).To(Equal("https://httpbin.org/anything"))
			})

			It("request to anything by option style", func() {
				var result demo.AnythingResponse

				_, err := client.Anything(
					bkapi.OptJsonResultProvider(),
					bkapi.OptSetRequestResult(&result),
				).Request()

				if err != nil {
					Skip(err.Error())
				}

				// it is also ok to set the result provider by option style
				Expect(result.URL).To(Equal("https://httpbin.org/anything"))
			})
		})

		Context("with request payload", func() {
			var (
				result *demo.AnythingResponse
			)

			BeforeEach(func() {
				result = new(demo.AnythingResponse)

				// you can set the common options to the client
				err := client.AddOperationOptions(
					bkapi.JsonBodyProvider(),
					bkapi.JsonResultProvider(),
					bkapi.OptSetRequestResult(result),
				)

				Expect(err).To(BeNil())
			})

			It("request to anything by chaining style", func() {
				_, err := client.Anything().
					SetQueryParams(map[string]string{
						"from": "query",
					}).
					SetBody(map[string]interface{}{
						"from": "body",
					}).
					SetHeaders(map[string]string{
						"X-Header": "my-header",
					}).
					Request()

				if err != nil {
					Skip(err.Error())
				}

				Expect(result.Args["from"]).To(Equal("query"))
				Expect(result.JSON["from"]).To(Equal("body"))
				Expect(result.Headers["X-Header"]).To(Equal("my-header"))
			})

			It("request to anything by option style", func() {
				_, err := client.Anything(
					bkapi.OptSetRequestQueryParams(map[string]string{
						"from": "query",
					}),
					bkapi.OptSetRequestBody(map[string]interface{}{
						"from": "body",
					}),
					bkapi.OptSetRequestHeaders(map[string]string{
						"X-Header": "my-header",
					}),
				).Request()

				if err != nil {
					Skip(err.Error())
				}

				Expect(result.Args["from"]).To(Equal("query"))
				Expect(result.JSON["from"]).To(Equal("body"))
				Expect(result.Headers["X-Header"]).To(Equal("my-header"))
			})
		})

		Context("with path params", func() {
			It("request to anything by chaining style", func() {
				response, err := client.StatusCode().
					SetPathParams(map[string]string{
						"code": "200",
					}).
					Request()

				if err != nil {
					Skip(err.Error())
				}

				Expect(response.Request.URL.Path).To(Equal("/status/200"))
			})

			It("request to anything by option style", func() {
				response, err := client.StatusCode(
					bkapi.OptSetRequestPathParams(map[string]string{
						"code": "200",
					}),
				).Request()

				if err != nil {
					Skip(err.Error())
				}

				Expect(response.Request.URL.Path).To(Equal("/status/200"))
			})
		})
	})

	Context("Example for error handling", func() {
		It("should handle 5XX", func() {
			response, err := client.StatusCode().
				SetPathParams(map[string]string{
					"code": "500",
				}).
				Request()

			if err != nil {
				Skip(err.Error())
			}

			Expect(response.StatusCode).To(Equal(500))
		})

		It("should handle decode error", func() {
			var result map[string]interface{}

			// 5xx response will sent non-json body
			_, err := client.StatusCode().
				SetResultProvider(bkapi.JsonResultProvider()).
				SetResult(&result).
				SetPathParams(map[string]string{
					"code": "500",
				}).
				Request()

			if err != nil {
				Skip(err.Error())
			}

			Expect(result).To(BeEmpty())
		})
	})
})
