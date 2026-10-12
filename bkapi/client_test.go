/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package bkapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gentleman "gopkg.in/h2non/gentleman.v2"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/bkapi"
)

func newClient(config bkapi.Config) *gentleman.Client {
	GinkgoHelper()
	if config.Logger == nil {
		config.Logger = slog.New(slog.DiscardHandler)
	}
	client, err := bkapi.New(config)
	Expect(err).To(BeNil())
	return client
}

// serve starts a server that replies with the status, headers and body, and returns its url.
func serve(status int, header http.Header, body string) string {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, values := range header {
			w.Header()[key] = values
		}
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	DeferCleanup(server.Close)
	return server.URL
}

var _ = Describe("Client", func() {
	ctx := context.Background()

	DescribeTable("should reject an invalid endpoint",
		func(endpoint string) {
			_, err := bkapi.New(bkapi.Config{Endpoint: endpoint})
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("without scheme", "example.com"),
		Entry("relative path", "/api/my-gateway/prod"),
		Entry("without host", "http://"),
	)

	Context("Request", func() {
		var (
			endpoint string
			request  *http.Request
			body     string
		)

		BeforeEach(func() {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, _ := io.ReadAll(r.Body)
				request, body = r.Clone(context.Background()), string(data)
				fmt.Fprint(w, `{"id": 1, "name": "prod"}`)
			}))
			DeferCleanup(server.Close)
			endpoint = server.URL
		})

		It("should send the authorization, tenant, path params and json body", func() {
			c := newClient(bkapi.Config{
				Endpoint:  endpoint + "/api/my-gateway/prod/",
				AppCode:   "my-app",
				AppSecret: "my-secret",
				TenantID:  "my-tenant",
			})
			res, err := c.Post().
				AddPath("/stages/:name/").
				Param("name", "prod").
				JSON(map[string]string{"description": "prod"}).
				Send()
			Expect(err).To(BeNil())

			var result struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}
			Expect(res.JSON(&result)).To(Succeed())
			Expect(result.ID).To(Equal(1))
			Expect(result.Name).To(Equal("prod"))

			Expect(request.Method).To(Equal(http.MethodPost))
			Expect(request.URL.Path).To(Equal("/api/my-gateway/prod/stages/prod/"))
			Expect(request.Header.Get("X-Bkapi-Authorization")).To(
				MatchJSON(`{"bk_app_code": "my-app", "bk_app_secret": "my-secret"}`))
			Expect(request.Header.Get("X-Bk-Tenant-Id")).To(Equal("my-tenant"))
			Expect(request.Header.Get("Content-Type")).To(Equal("application/json"))
			Expect(body).To(MatchJSON(`{"description": "prod"}`))
		})

		It("should send the access token without the tenant", func() {
			c := newClient(bkapi.Config{Endpoint: endpoint, AppCode: "my-app", AccessToken: "my-token"})
			_, err := c.Get().Send()
			Expect(err).To(BeNil())

			Expect(request.Header.Get("X-Bkapi-Authorization")).To(
				MatchJSON(`{"bk_app_code": "my-app", "access_token": "my-token"}`))
			Expect(request.Header.Get("X-Bk-Tenant-Id")).To(BeEmpty())
		})

		It("should send the query", func() {
			c := newClient(bkapi.Config{Endpoint: endpoint})
			_, err := c.Get().AddPath("/versions/").SetQuery("version", "1.0.0+prod").Send()
			Expect(err).To(BeNil())

			Expect(request.URL.Query().Get("version")).To(Equal("1.0.0+prod"))
		})
	})

	DescribeTable("should preserve the error of a non-2xx response",
		func(status int, header http.Header, body, wantData, requestID string) {
			c := newClient(bkapi.Config{Endpoint: serve(status, header, body)})
			_, err := c.Get().Send()

			var apiErr *bkapi.Error
			Expect(errors.As(err, &apiErr)).To(BeTrue())
			var data any
			Expect(json.Unmarshal([]byte(wantData), &data)).To(Succeed())
			Expect(*apiErr).To(Equal(bkapi.Error{StatusCode: status, Data: data, RequestID: requestID}))
		},
		Entry("from the response body", 400,
			http.Header{"X-Bkapi-Request-Id": {"request-1"}},
			`{"error": {"code": "INVALID_ARGUMENT", "message": "invalid version", "details": [], "data": {}}}`,
			`{"code": "INVALID_ARGUMENT", "message": "invalid version", "details": [], "data": {}}`, "request-1"),
		Entry("with fields defined by the API", 403,
			http.Header{"X-Request-Id": {"request-3"}},
			`{"error": {"code": 403, "message": ["denied"], "system": "bk-job",
    "details": {"action": "read"}, "data": {"system": "bk_job"}, "retryable": false}}`,
			`{"code": 403, "message": ["denied"], "system": "bk-job",
    "details": {"action": "read"}, "data": {"system": "bk_job"}, "retryable": false}`, "request-3"),
		Entry("without code or message", 400, nil,
			`{"error": {"reason": "invalid name", "field": "name"}}`,
			`{"reason": "invalid name", "field": "name"}`, ""),
		Entry("with an empty object", 400, nil, `{"error": {}}`, `{}`, ""),
		Entry("with a null error", 400, nil, `{"error": null}`, `null`, ""),
		Entry("with a string error", 400, nil, `{"error": "invalid name"}`, `"invalid name"`, ""),
		Entry("with a list error", 400, nil, `{"error": [{"field": "name"}]}`, `[{"field": "name"}]`, ""),
		Entry("from the legacy body of the gateway", 403,
			http.Header{"X-Bkapi-Request-Id": {"request-2"}, "X-Bkapi-Error-Code": {"1640301"}},
			`{"code": 1640301, "code_name": "APP_NO_PERMISSION", "message": "App has no permission"}`,
			`{"code": 1640301, "code_name": "APP_NO_PERMISSION", "message": "App has no permission"}`, "request-2"),
		Entry("from the legacy body with a numeric code", 400, nil,
			`{"result": false, "code": 40000, "message": "invalid params", "data": null}`,
			`{"result": false, "code": 40000, "message": "invalid params", "data": null}`, ""),
		Entry("from the headers of the gateway", 401,
			http.Header{"X-Bkapi-Error-Code": {"1640001"}, "X-Bkapi-Error-Message": {"App authentication failed"}},
			``, `{"code": "1640001", "message": "App authentication failed"}`, ""),
		Entry("from the status", 502, nil, `<html>502 Bad Gateway</html>`,
			`{"code": "", "message": "Bad Gateway"}`, ""),
		Entry("of a redirect", 302, nil, ``, `{"code": "", "message": "Found"}`, ""),
	)

	Context("Server not responding", func() {
		var endpoint string

		BeforeEach(func() {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				<-r.Context().Done()
			}))
			DeferCleanup(server.Close)
			endpoint = server.URL
		})

		It("should stop when the context is canceled", func() {
			canceled, cancel := context.WithCancel(ctx)
			cancel()

			c := newClient(bkapi.Config{Endpoint: endpoint})
			_, err := c.Get().Use(bkapi.WithContext(canceled)).Send()
			Expect(err).To(MatchError(context.Canceled))
		})

		It("should stop when the request times out", func() {
			c := newClient(bkapi.Config{Endpoint: endpoint, Timeout: 50 * time.Millisecond})
			_, err := c.Get().Send()
			Expect(err).To(MatchError(context.DeadlineExceeded))
		})
	})

	Context("Slow error response", func() {
		var (
			endpoint string
			started  chan struct{} // closed when the server starts to send the body of /slow
		)

		BeforeEach(func() {
			started = make(chan struct{})
			release := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/slow" {
					return
				}
				w.WriteHeader(http.StatusServiceUnavailable)
				w.(http.Flusher).Flush()
				close(started)
				select {
				case <-release:
				case <-r.Context().Done():
				}
			}))
			// the cleanups run in reverse order, so the handler is released before the server is closed
			DeferCleanup(server.Close)
			DeferCleanup(func() { close(release) })
			endpoint = server.URL
		})

		It("should not block the other requests of the client", func() {
			c := newClient(bkapi.Config{Endpoint: endpoint})
			go func() { _, _ = c.Get().AddPath("/slow").Send() }()
			Eventually(started).Should(BeClosed())

			done := make(chan error)
			go func() {
				_, err := c.Get().AddPath("/fast").Send()
				done <- err
			}()
			Eventually(done, time.Second).Should(Receive(BeNil()))
		})

		It("should stop reading the body when the request times out", func() {
			c := newClient(bkapi.Config{Endpoint: endpoint, Timeout: 50 * time.Millisecond})
			_, err := c.Get().AddPath("/slow").Send()
			Expect(err).To(MatchError(context.DeadlineExceeded))
		})
	})

	It("should send the requests through the transport", func() {
		var used bool
		transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
			used = true
			return http.DefaultTransport.RoundTrip(r)
		})

		c := newClient(bkapi.Config{Endpoint: serve(200, nil, `[]`), Transport: transport})
		_, err := c.Get().Send()
		Expect(err).To(BeNil())
		Expect(used).To(BeTrue())
	})

	DescribeTable("should log the requests without the secrets",
		func(status int, level string) {
			var logs bytes.Buffer
			c := newClient(bkapi.Config{
				Endpoint:  serve(status, http.Header{"X-Bkapi-Request-Id": {"request-1"}}, `{}`),
				AppSecret: "my-secret",
				Logger:    slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})),
			})
			_, _ = c.Post().
				AddPath("/gateways/my-gateway/").
				JSON(map[string]string{"description": "my-description"}).
				Send()

			Expect(logs.String()).To(And(
				ContainSubstring("level="+level),
				ContainSubstring("method=POST"),
				ContainSubstring("path=/gateways/my-gateway/"),
				ContainSubstring(fmt.Sprint("status=", status)),
				ContainSubstring("request_id=request-1"),
			))
			Expect(logs.String()).NotTo(Or(ContainSubstring("my-secret"), ContainSubstring("my-description")))
			Expect(bytes.Count(logs.Bytes(), []byte("\n"))).To(Equal(1))
		},
		Entry("successful", 200, "DEBUG"),
		Entry("client error", 400, "WARN"),
		Entry("server error", 500, "ERROR"),
	)

	It("should log the error of the gateway", func() {
		var logs bytes.Buffer
		c := newClient(bkapi.Config{
			Endpoint: serve(403, http.Header{"X-Bkapi-Error-Code": {"APP_NO_PERMISSION"}}, ``),
			Logger:   slog.New(slog.NewTextHandler(&logs, nil)),
		})
		_, err := c.Get().Send()
		Expect(err).To(HaveOccurred())
		Expect(logs.String()).To(ContainSubstring("APP_NO_PERMISSION"))
	})

	It("should log the complete error object and the traceparent", func() {
		var logs bytes.Buffer
		c := newClient(bkapi.Config{
			Endpoint: serve(400, nil, `{"error": {"reason": "invalid", "details": {"name": "required"}}}`),
			Logger:   slog.New(slog.NewJSONHandler(&logs, nil)),
		})
		_, err := c.Get().SetHeader("traceparent", "00-trace-span-01").Send()

		Expect(err).To(MatchError(ContainSubstring("required")))
		var entry map[string]any
		Expect(json.Unmarshal(logs.Bytes(), &entry)).To(Succeed())
		Expect(entry["error"]).To(Equal(map[string]any{
			"reason": "invalid", "details": map[string]any{"name": "required"},
		}))
		Expect(entry["traceparent"]).To(Equal("00-trace-span-01"))
	})

	It("should not log the errors without a response", func() {
		var logs bytes.Buffer
		c := newClient(bkapi.Config{
			Endpoint: "http://example.com",
			Logger:   slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})),
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, errors.New("connection refused")
			}),
		})
		_, err := c.Get().SetQuery("token", "my-token").Send()
		Expect(err).To(MatchError(ContainSubstring("connection refused")))
		Expect(logs.String()).To(BeEmpty())
	})

	It("should be safe for concurrent use", func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `%q`, r.URL.Path)
		}))
		DeferCleanup(server.Close)

		c := newClient(bkapi.Config{Endpoint: server.URL})
		var wg sync.WaitGroup
		for i := range 20 {
			wg.Go(func() {
				defer GinkgoRecover()
				path := fmt.Sprint("/items/", i, "/")
				res, err := c.Get().AddPath(path).Send()
				Expect(err).To(BeNil())

				var result string
				Expect(res.JSON(&result)).To(Succeed())
				Expect(result).To(Equal(path))
			})
		}
		wg.Wait()
	})
})

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
