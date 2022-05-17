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

package internal_test

import (
	"fmt"
	"net/http"

	"github.com/TencentBlueKing/gopkg/logging"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/transport"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal/mock"
)

var _ = Describe("Client", func() {
	var (
		ctrl                    *gomock.Controller
		client                  *internal.BkApiClient
		gentlemanClient         *gentleman.Client
		operation               *mock.MockOperation
		request                 *gentleman.Request
		response                *http.Response
		mockTransport           *mock.MockRoundTripper
		clientConfig            *mock.MockClientConfig
		operationConfig         *mock.MockOperationConfig
		operationConfigProvider *mock.MockOperationConfigProvider
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockTransport = mock.NewMockRoundTripper(ctrl)
		operation = mock.NewMockOperation(ctrl)

		gentlemanClient = gentleman.New()
		gentlemanClient.Use(transport.Set(mockTransport))

		clientConfig = mock.NewMockClientConfig(ctrl)
		clientConfig.EXPECT().GetUrl().Return("").AnyTimes()
		clientConfig.EXPECT().GetAuthorizationHeaders().Return(nil).AnyTimes()
		clientConfig.EXPECT().GetLogger().Return(logging.GetLogger("")).AnyTimes()

		client = internal.NewBkApiClient(
			"testing", gentlemanClient,
			func(name string, req *gentleman.Request) define.Operation {
				operation.EXPECT().Name().Return(name).AnyTimes()

				request = req
				return operation
			},
			clientConfig,
		)

		operationConfig = mock.NewMockOperationConfig(ctrl)
		operationConfig.EXPECT().GetName().Return("").AnyTimes()
		operationConfig.EXPECT().GetPath().Return("").AnyTimes()
		operationConfig.EXPECT().GetMethod().Return("").AnyTimes()

		operationConfigProvider = mock.NewMockOperationConfigProvider(ctrl)
		operationConfigProvider.EXPECT().ProvideConfig().Return(operationConfig).AnyTimes()

		response = &http.Response{
			Header: http.Header{},
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	mockTransportRoundTrip := func() {
		mockTransport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			response.Request = req

			return response, nil
		})
	}

	Context("BkApiClient", func() {
		AfterEach(func() {
			ctrl.Finish()
		})

		It("should fail on apply", func() {
			option := mock.NewMockBkApiClientOption(ctrl)

			option.EXPECT().ApplyToClient(gomock.Any()).Return(errors.New("test"))
			Expect(client.Apply(option)).NotTo(BeNil())
		})

		It("should apply an option", func() {
			option := mock.NewMockBkApiClientOption(ctrl)
			option.EXPECT().ApplyToClient(client).Return(nil)

			Expect(client.Apply(option)).To(BeNil())
		})

		It("should new an operation", func() {
			operationConfig := mock.NewMockOperationConfig(ctrl)
			operationConfigProvider := mock.NewMockOperationConfigProvider(ctrl)
			operationConfigProvider.EXPECT().ProvideConfig().Return(operationConfig).AnyTimes()

			operationConfig.EXPECT().GetName().Return("test").AnyTimes()
			operationConfig.EXPECT().GetPath().Return("/test").AnyTimes()
			operationConfig.EXPECT().GetMethod().Return("POST").AnyTimes()

			operation.EXPECT().Apply(gomock.Any()).AnyTimes()

			op := client.NewOperation(operationConfigProvider)
			Expect(op).To(Equal(operation))

			mockTransportRoundTrip()
			// make a request to apply the middlewares
			_, err := request.Send()
			Expect(err).To(BeNil())

			Expect(request.Context.Request.Method).To(Equal("POST"))
			Expect(request.Context.Request.URL.Path).To(Equal("/test"))
		})

		It("should new an operation with option", func() {
			option := mock.NewMockOperationOption(ctrl)
			operation.EXPECT().Apply(option)

			op := client.NewOperation(operationConfigProvider, option)
			Expect(op).NotTo(BeNil())
		})

		It("should new an operation with common operation option", func() {
			option := mock.NewMockOperationOption(ctrl)
			operation.EXPECT().Apply(option)
			Expect(client.AddOperationOptions(option)).To(BeNil())

			op := client.NewOperation(operationConfigProvider)
			Expect(op).NotTo(BeNil())
		})

		It("should generate operation name by config.Name", func() {
			operationConfig := mock.NewMockOperationConfig(ctrl)
			operationConfigProvider := mock.NewMockOperationConfigProvider(ctrl)
			operationConfigProvider.EXPECT().ProvideConfig().Return(operationConfig).AnyTimes()

			operationConfig.EXPECT().GetName().Return("operation").AnyTimes()
			operationConfig.EXPECT().GetMethod().Return("GET").AnyTimes()
			operationConfig.EXPECT().GetPath().Return("/test").AnyTimes()

			operation := client.NewOperation(operationConfigProvider)

			Expect(operation.Name()).To(Equal("testing.operation"))
		})

		It("should generate anonymous operation name", func() {
			operationConfig := mock.NewMockOperationConfig(ctrl)
			operationConfigProvider := mock.NewMockOperationConfigProvider(ctrl)
			operationConfigProvider.EXPECT().ProvideConfig().Return(operationConfig).AnyTimes()

			operationConfig.EXPECT().GetName().Return("").AnyTimes()
			operationConfig.EXPECT().GetMethod().Return("GET").AnyTimes()
			operationConfig.EXPECT().GetPath().Return("/test").AnyTimes()

			operation := client.NewOperation(operationConfigProvider)

			Expect(operation.Name()).To(Equal("testing(GET /test)"))
		})

		It("should set the user agent", func() {
			op := client.NewOperation(operationConfigProvider)
			Expect(op).To(Equal(operation))

			mockTransportRoundTrip()

			_, err := request.Send()
			Expect(err).To(BeNil())

			request := request.Context.Request
			Expect(request.Header.Get("User-Agent")).To(Equal(fmt.Sprintf("%s/%s", define.UserAgent, define.Version)))
		})

		Context("Logger", func() {
			var (
				logger *mock.MockLogger
			)

			BeforeEach(func() {
				logger = mock.NewMockLogger(ctrl)

				clientConfig = mock.NewMockClientConfig(ctrl)
				clientConfig.EXPECT().GetUrl().Return("").AnyTimes()
				clientConfig.EXPECT().GetAuthorizationHeaders().Return(nil).AnyTimes()
				clientConfig.EXPECT().GetLogger().Return(logger).AnyTimes()

				client = internal.NewBkApiClient(
					"testing", gentlemanClient,
					func(name string, req *gentleman.Request) define.Operation {
						operation.EXPECT().Name().Return(name).AnyTimes()

						request = req
						return operation
					},
					clientConfig,
				)
			})

			It("should log 2xx response details", func() {
				logger.EXPECT().
					Debug(gomock.Any(), gomock.Any()).
					DoAndReturn(func(msg string, fields ...map[string]interface{}) {
						Expect(msg).To(Equal("request success"))
						Expect(fields).To(HaveLen(1))

						values := fields[0]
						Expect(values["status_code"]).To(Equal(200))
						Expect(values["operation"]).To(Equal(operation))
						Expect(values["bkapi_request_id"]).To(
							Equal(response.Header.Get("X-Bkapi-Request-Id")),
						)
					})

				op := client.NewOperation(operationConfigProvider)
				Expect(op).To(Equal(operation))

				response.StatusCode = 200
				response.Header.Set("X-Bkapi-Request-Id", "12345")
				mockTransportRoundTrip()

				_, err := request.Send()
				Expect(err).To(BeNil())
			})

			It("should log 4xx response details", func() {
				logger.EXPECT().
					Warn(gomock.Any(), gomock.Any()).
					DoAndReturn(func(msg string, fields ...map[string]interface{}) {
						Expect(msg).To(Equal("request error caused by client"))
						Expect(fields).To(HaveLen(1))

						values := fields[0]
						Expect(values["status_code"]).To(Equal(403))
						Expect(values["operation"]).To(Equal(operation))
						Expect(values["bkapi_request_id"]).To(
							Equal(response.Header.Get("X-Bkapi-Request-Id")),
						)
						Expect(values["bkapi_error_code"]).To(
							Equal(response.Header.Get("X-Bkapi-Error-Code")),
						)
						Expect(values["bkapi_error_message"]).To(
							Equal(response.Header.Get("X-Bkapi-Error-Message")),
						)
					})

				op := client.NewOperation(operationConfigProvider)
				Expect(op).To(Equal(operation))

				response.StatusCode = 403
				response.Header.Set("X-Bkapi-Request-Id", "12345")
				response.Header.Set("X-Bkapi-Error-Code", "0")
				response.Header.Set("X-Bkapi-Error-Message", "test error message")
				mockTransportRoundTrip()

				_, err := request.Send()
				Expect(err).To(BeNil())
			})
			It("should log 5xx response details", func() {
				logger.EXPECT().
					Error(gomock.Any(), gomock.Any()).
					DoAndReturn(func(msg string, fields ...map[string]interface{}) {
						Expect(msg).To(Equal("request error caused by server"))
						Expect(fields).To(HaveLen(1))

						values := fields[0]
						Expect(values["status_code"]).To(Equal(502))
						Expect(values["operation"]).To(Equal(operation))
						Expect(values["bkapi_request_id"]).To(
							Equal(response.Header.Get("X-Bkapi-Request-Id")),
						)
						Expect(values["bkapi_error_code"]).To(
							Equal(response.Header.Get("X-Bkapi-Error-Code")),
						)
						Expect(values["bkapi_error_message"]).To(
							Equal(response.Header.Get("X-Bkapi-Error-Message")),
						)
					})

				op := client.NewOperation(operationConfigProvider)
				Expect(op).To(Equal(operation))

				response.StatusCode = 502
				response.Header.Set("X-Bkapi-Request-Id", "12345")
				response.Header.Set("X-Bkapi-Error-Code", "0")
				response.Header.Set("X-Bkapi-Error-Message", "test error message")
				mockTransportRoundTrip()

				_, err := request.Send()
				Expect(err).To(BeNil())
			})
		})
	})

	Context("BkApiClientOption", func() {
		It("should fail when the type is not supported", func() {
			var client mock.MockBkApiClient

			opt := internal.NewBkApiClientOption(nil)
			err := opt.ApplyToClient(&client)

			Expect(errors.Cause(err)).To(Equal(define.ErrTypeNotMatch))
		})

		It("should apply function", func() {
			var client internal.BkApiClient
			err := fmt.Errorf("testing")

			opt := internal.NewBkApiClientOption(func(c *internal.BkApiClient) error {
				Expect(c).To(Equal(&client))
				return err
			})

			Expect(opt.ApplyToClient(&client)).To(Equal(err))
		})
	})
})
