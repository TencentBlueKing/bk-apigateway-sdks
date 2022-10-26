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
	"strconv"
	"sync"
	"time"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/h2non/gentleman.v2/context"
)

// PrometheusOptions for common prometheus metrics options
type PrometheusOptions struct {
	Namespace       string
	Subsystem       string
	ConstLabels     prometheus.Labels
	DurationBuckets []float64
	BytesBuckets    []float64
	Registerer      prometheus.Registerer
}

type bkapiCollector struct {
	*internal.OperationOption
	metricRequestsDurationSeconds *prometheus.HistogramVec
	metricRequestsBodyBytes       *prometheus.HistogramVec
	metricResponsesBodyBytes      *prometheus.HistogramVec
	metricResponsesTotal          *prometheus.CounterVec
	metricResponsesFailuresTotal  *prometheus.CounterVec
}

func (c *bkapiCollector) init(opt PrometheusOptions) {
	c.OperationOption = internal.NewOperationOption(c.collectMetrics)

	registerer := opt.Registerer

	c.metricRequestsDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   opt.Namespace,
			Subsystem:   opt.Subsystem,
			ConstLabels: opt.ConstLabels,
			Buckets:     opt.DurationBuckets,
			Name:        "bkapi_requests_duration_seconds",
			Help:        "Histogram of requests duration by operation, method",
		}, []string{"operation", "method"},
	)
	registerer.MustRegister(c.metricRequestsDurationSeconds)

	c.metricRequestsBodyBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   opt.Namespace,
			Subsystem:   opt.Subsystem,
			ConstLabels: opt.ConstLabels,
			Buckets:     opt.BytesBuckets,
			Name:        "bkapi_requests_body_bytes",
			Help:        "Histogram of requests body bytes by operation, method",
		}, []string{"operation", "method"},
	)
	registerer.MustRegister(c.metricRequestsBodyBytes)

	c.metricResponsesBodyBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   opt.Namespace,
			Subsystem:   opt.Subsystem,
			ConstLabels: opt.ConstLabels,
			Buckets:     opt.BytesBuckets,
			Name:        "bkapi_responses_body_bytes",
			Help:        "Histogram of responses body bytes by operation, method",
		}, []string{"operation", "method"},
	)
	registerer.MustRegister(c.metricResponsesBodyBytes)

	c.metricResponsesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   opt.Namespace,
			Subsystem:   opt.Subsystem,
			ConstLabels: opt.ConstLabels,
			Name:        "bkapi_responses_total",
			Help:        "Count of responses by operation, method, status",
		}, []string{"operation", "method", "status"},
	)
	registerer.MustRegister(c.metricResponsesTotal)

	c.metricResponsesFailuresTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   opt.Namespace,
			Subsystem:   opt.Subsystem,
			ConstLabels: opt.ConstLabels,
			Name:        "bkapi_failures_total",
			Help:        "Count of failures by operation, method",
		}, []string{"operation", "method", "error"},
	)
	registerer.MustRegister(c.metricResponsesFailuresTotal)
}

func (c *bkapiCollector) collectMetrics(operation *internal.Operation) error {
	var (
		requestStart time.Time
		name         = operation.FullName()
		request      = internal.GetOperationRawRequest(operation)
	)

	request.UseHandler("before dial", func(ctx *context.Context, h context.Handler) {
		requestStart = time.Now()
		h.Next(ctx)
	})

	request.UseHandler("after dial", func(ctx *context.Context, h context.Handler) {
		requestEnd := time.Now()
		defer h.Next(ctx)

		method := ctx.Request.Method
		status := strconv.Itoa(ctx.Response.StatusCode)
		c.metricResponsesTotal.WithLabelValues(name, method, status).Inc()

		if !requestStart.IsZero() {
			c.metricRequestsDurationSeconds.WithLabelValues(name, method).
				Observe(requestEnd.Sub(requestStart).Seconds())
		}

		requestContentLength, err := strconv.ParseFloat(ctx.Request.Header.Get("Content-Length"), 64)
		if err == nil {
			c.metricRequestsBodyBytes.WithLabelValues(name, method).Observe(requestContentLength)
		}

		responseContentLength, err := strconv.ParseFloat(ctx.Response.Header.Get("Content-Length"), 64)
		if err == nil {
			c.metricResponsesBodyBytes.WithLabelValues(name, method).Observe(responseContentLength)
		}
	})

	request.UseHandler("error", func(ctx *context.Context, h context.Handler) {
		defer h.Next(ctx)

		cause := define.ErrorCause(ctx.Error)
		if cause == nil {
			return
		}

		c.metricResponsesFailuresTotal.WithLabelValues(name, ctx.Request.Method, cause.Error()).Inc()
	})

	return nil
}

func initCollector(collector *bkapiCollector, opt PrometheusOptions) {
	if opt.DurationBuckets == nil {
		opt.DurationBuckets = []float64{
			1,
			5,
			10,
			50,
			100,
			500,
			1024,       // 1k
			5120,       // 5k
			10240,      // 10k
			51200,      // 50k
			102400,     // 100k
			512000,     // 500k
			1048576,    // 1m
			5242880,    // 5m
			10485760,   // 10m
			52428800,   // 50m
			104857600,  // 100m
			524288000,  // 500m
			1073741824, // 1g
		}
	}

	if opt.BytesBuckets == nil {
		opt.BytesBuckets = []float64{
			0.1,
			0.25,
			0.5,
			0.75,
			1.0,
			2.5,
			5.0,
			7.5,
			10.0,
			25.0,
			50.0,
			75.0,
			100.0,
			250.0,
			500.0,
			750.0,
		}
	}

	if opt.Registerer == nil {
		opt.Registerer = prometheus.DefaultRegisterer
	}

	collector.init(opt)
}

var (
	initOnce             sync.Once
	globalBkapiCollector bkapiCollector
)

// Enable prometheus metrics
func Enable(opt PrometheusOptions) (ok bool) {
	initOnce.Do(func() {
		initCollector(&globalBkapiCollector, opt)
		bkapi.RegisterGlobalBkapiClientOption(globalBkapiCollector)

		ok = true
	})

	return ok
}
