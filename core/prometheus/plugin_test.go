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

package prometheus

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

var _ = Describe("Plugin", func() {
	var (
		ctrl            *gomock.Controller
		mockTransport   *mock.MockRoundTripper
		response        *http.Response
		requestError    error
		registry        *prometheus.Registry
		collector       *bkapiCollector
		client          define.BkApiClient
		apiName         = "testing"
		operationName   string
		clientConfig    bkapi.ClientConfig
		operationConfig bkapi.OperationConfig
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockTransport = mock.NewMockRoundTripper(ctrl)
		requestError = nil
		response = &http.Response{
			StatusCode: 200,
			Header:     http.Header{},
		}

		collector = &bkapiCollector{}
		registry = prometheus.NewRegistry()
		initCollector(collector, PrometheusOptions{
			Registerer: registry,
		})

		operationConfig = bkapi.OperationConfig{
			Name:   "demo",
			Method: "POST",
			Path:   "/api/v1/demo",
		}

		clientConfig = bkapi.ClientConfig{
			Endpoint: "http://example.com",
		}
		var err error
		client, err = bkapi.NewBkApiClient(apiName, clientConfig, collector, bkapi.OptTransport(mockTransport))
		Expect(err).To(BeNil())

		operationName = fmt.Sprintf("%s.api.%s", client.Name(), operationConfig.Name)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	mockRequest := func() {
		mockTransport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			response.Request = req

			return response, requestError
		})
	}

	gatherMetric := func(metricName string, labels map[string]string) *io_prometheus_client.Metric {
		results, err := registry.Gather()
		Expect(err).To(BeNil())

		for _, result := range results {
			if *result.Name != metricName {
				continue
			}

		outer:
			for _, metric := range result.Metric {
				for _, label := range metric.Label {
					value, ok := labels[*label.Name]
					if !ok {
						continue
					}

					if value != *label.Value {
						break outer
					}
				}

				return metric
			}
		}

		return nil
	}

	Context("bkapi_requests_duration_seconds", func() {
		It("should record when success", func() {
			mockRequest()
			_, err := client.NewOperation(operationConfig).Request()
			Expect(err).To(BeNil())

			metric := gatherMetric("bkapi_requests_duration_seconds", map[string]string{
				"operation": operationName,
				"method":    operationConfig.Method,
			})
			Expect(metric).NotTo(BeNil())
		})

		It("should not record when error", func() {
			requestError = fmt.Errorf("testing")

			mockRequest()
			_, err := client.NewOperation(operationConfig).Request()
			Expect(err).NotTo(BeNil())

			metric := gatherMetric("bkapi_requests_duration_seconds", map[string]string{
				"operation": operationName,
				"method":    operationConfig.Method,
			})
			Expect(metric).To(BeNil())
		})
	})

	Context("bkapi_responses_total", func() {
		DescribeTable("should record when success", func(statusCode int) {
			response.StatusCode = statusCode
			mockRequest()
			response, err := client.NewOperation(operationConfig).Request()
			Expect(err).To(BeNil())

			metric := gatherMetric("bkapi_responses_total", map[string]string{
				"operation": operationName,
				"method":    operationConfig.Method,
				"status":    strconv.Itoa(response.StatusCode),
			})
			Expect(metric).NotTo(BeNil())
		},
			Entry("200", 200),
			Entry("400", 400),
			Entry("500", 500),
		)

		It("should not record when error", func() {
			requestError = fmt.Errorf("testing")

			mockRequest()
			_, err := client.NewOperation(operationConfig).Request()
			Expect(err).NotTo(BeNil())

			metric := gatherMetric("bkapi_responses_total", map[string]string{
				"operation": operationName,
				"method":    operationConfig.Method,
				"status":    strconv.Itoa(response.StatusCode),
			})
			Expect(metric).To(BeNil())
		})
	})

	Context("bkapi_failures_total", func() {
		DescribeTable("should not record when success", func(statusCode int) {
			response.StatusCode = statusCode
			mockRequest()
			_, err := client.NewOperation(operationConfig).Request()
			Expect(err).To(BeNil())

			metric := gatherMetric("bkapi_failures_total", map[string]string{
				"operation": operationName,
				"method":    operationConfig.Method,
			})
			Expect(metric).To(BeNil())
		},
			Entry("200", 200),
			Entry("400", 400),
			Entry("500", 500),
		)

		It("should record when error", func() {
			requestError = fmt.Errorf("testing")

			mockRequest()
			_, err := client.NewOperation(operationConfig).Request()
			Expect(err).NotTo(BeNil())

			metric := gatherMetric("bkapi_failures_total", map[string]string{
				"operation": operationName,
				"method":    operationConfig.Method,
			})
			Expect(metric).NotTo(BeNil())
		})
	})

	Context("bkapi_requests_body_bytes", func() {
		It("should record when contain Content-Length header", func() {
			mockRequest()
			_, err := client.NewOperation(operationConfig).
				SetHeaders(map[string]string{"Content-Length": "1024"}).
				Request()
			Expect(err).To(BeNil())

			metric := gatherMetric("bkapi_requests_body_bytes", map[string]string{
				"operation": operationName,
				"method":    operationConfig.Method,
			})
			Expect(metric).NotTo(BeNil())
		})
	})

	Context("bkapi_responses_body_bytes", func() {
		It("should record when contain Content-Length header", func() {
			response.Header.Add("Content-Length", "1024")

			mockRequest()
			_, err := client.NewOperation(operationConfig).Request()
			Expect(err).To(BeNil())

			metric := gatherMetric("bkapi_responses_body_bytes", map[string]string{
				"operation": operationName,
				"method":    operationConfig.Method,
			})
			Expect(metric).NotTo(BeNil())
		})
	})
})
