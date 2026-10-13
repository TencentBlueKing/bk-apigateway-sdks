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

package manager_test

import (
	"context"
	"encoding/json"
	"io"
	"maps"
	"net/http"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gock "gopkg.in/h2non/gock.v1"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/bkapi"
	manager "github.com/TencentBlueKing/bk-apigateway-sdks/v2/manager"
)

var _ = Describe("Manager", func() {
	var (
		ctx      = context.Background()
		endpoint = "http://example.com"
		mgr      *manager.Manager
		bodies   []map[string]any
	)

	captureBody := func(req *http.Request, _ *gock.Request) (bool, error) {
		if req.Body == nil {
			return true, nil
		}
		content, err := io.ReadAll(req.Body)
		if err != nil {
			return false, err
		}
		body := map[string]any{}
		if len(content) > 0 {
			if err := json.Unmarshal(content, &body); err != nil {
				return false, err
			}
		}
		bodies = append(bodies, body)
		return true, nil
	}

	newManager := func(definition map[string]any) *manager.Manager {
		m, err := manager.NewManager(
			"testing",
			bkapi.Config{Endpoint: endpoint, AppCode: "my-app", Transport: gock.NewTransport()},
			manager.NewDefinition(definition),
		)
		Expect(err).To(BeNil())
		return m
	}

	BeforeEach(func() {
		bodies = nil
		mgr = newManager(map[string]any{
			"stages": []any{
				map[string]any{"name": "prod"},
			},
			"grant_permissions": []any{
				map[string]any{"target_app_code": "app1", "grant_dimension": "gateway"},
				map[string]any{
					"target_app_code": "app2",
					"grant_dimension": "resource",
					"resource_names":  []any{"r1"},
				},
			},
			"apply_permissions": []any{
				map[string]any{"gateway_name": "other-gateway"},
			},
			"related_apps": []any{"app1"},
		})
	})

	AfterEach(func() {
		done := gock.IsDone()
		gock.Off()
		Expect(done).To(BeTrue())
	})

	It("should return the data of a v2 response", func() {
		gock.New(endpoint).
			Get("/api/v2/sync/gateways/testing/resource_versions/latest/").
			Reply(200).
			JSON(map[string]any{"data": map[string]any{"version": "1.0.0+prod"}})

		result, err := mgr.GetLatestResourceVersion(ctx)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(map[string]any{"version": "1.0.0+prod"}))
	})

	DescribeTable("should reject an empty successful JSON response",
		func(body string) {
			gock.New(endpoint).Get("/api/v2/sync/gateways/testing/resource_versions/latest/").
				Reply(200).BodyString(body)
			gock.New(endpoint).Get("/api/v2/sync/gateways/testing/resource_versions/").
				Reply(200).BodyString(body)

			_, err := mgr.GetLatestResourceVersion(ctx)
			Expect(err).To(MatchError(ContainSubstring("failed to decode response")))
			_, err = mgr.ResourceVersionExists(ctx, "1.0.0")
			Expect(err).To(MatchError(ContainSubstring("failed to decode response")))
		},
		Entry("empty body", ""),
		Entry("whitespace only", " \n\t"),
	)

	It("should accept an explicit null data value", func() {
		gock.New(endpoint).Get("/api/v2/sync/gateways/testing/resource_versions/latest/").
			Reply(200).BodyString(`{"data":null}`)
		result, err := mgr.GetLatestResourceVersion(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeNil())
	})

	It("should return the error of a v2 error response", func() {
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/resource_versions/").
			Reply(400).
			JSON(map[string]any{
				"error": map[string]any{"code": "INVALID_ARGUMENT", "message": "invalid version"},
			})

		_, err := mgr.CreateResourceVersion(ctx, "1.0.0+prod+20261008", "")
		Expect(err).To(BeAssignableToTypeOf(&bkapi.Error{}))
		Expect(err.Error()).To(ContainSubstring("INVALID_ARGUMENT"))
		Expect(err.Error()).To(ContainSubstring("invalid version"))
	})

	It("should return an error for a non-2xx response without an error body", func() {
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/resource_versions/release/").
			Reply(500).
			JSON(map[string]any{"data": nil})

		_, err := mgr.Release(ctx, "1.0.0+prod", "")
		Expect(err).To(BeAssignableToTypeOf(&bkapi.Error{}))
		Expect(err.Error()).To(ContainSubstring("500"))
	})

	It("should sync the stage mcp servers whose response data is a list", func() {
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/stages/prod/mcp-servers/").
			Reply(200).
			JSON(map[string]any{
				"data": []any{
					map[string]any{"id": 1, "name": "testing-prod-mcp", "action": "created"},
				},
			})

		result, err := newManager(map[string]any{
			"stages": []any{
				map[string]any{"name": "prod", "mcp_servers": []any{map[string]any{"name": "mcp"}}},
				// the stage without mcp_servers is skipped, otherwise gock fails the unexpected request
				map[string]any{"name": "test"},
			},
		}).SyncStageMcpConfig(ctx)
		Expect(err).To(BeNil())
		Expect(result).To(HaveLen(1))
		Expect(result["prod"]).To(HaveLen(1))
	})

	It("should check whether the resource version exists", func() {
		gock.New(endpoint).
			Get("/api/v2/sync/gateways/testing/resource_versions/").
			MatchParam("version", `^1\.0\.0\+prod$`).
			Reply(200).
			JSON(map[string]any{"data": map[string]any{"count": 1, "results": []any{}}})
		gock.New(endpoint).
			Get("/api/v2/sync/gateways/testing/resource_versions/").
			MatchParam("version", `^1\.0\.1\+prod$`).
			Reply(200).
			JSON(map[string]any{"data": map[string]any{"count": 0, "results": []any{}}})

		exists, err := mgr.ResourceVersionExists(ctx, "1.0.0+prod")
		Expect(err).To(BeNil())
		Expect(exists).To(BeTrue())

		exists, err = mgr.ResourceVersionExists(ctx, "1.0.1+prod")
		Expect(err).To(BeNil())
		Expect(exists).To(BeFalse())
	})

	It("should grant permissions with the v2 fields", func() {
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/permissions/grant/").
			Times(2).
			AddMatcher(captureBody).
			Reply(201).
			JSON(map[string]any{"data": nil})

		Expect(mgr.GrantPermissions(ctx)).To(Succeed())
		Expect(bodies).To(Equal([]map[string]any{
			{"target_app_code": "app1", "grant_dimension": "gateway"},
			{"target_app_code": "app2", "grant_dimension": "resource", "resource_names": []any{"r1"}},
		}))
	})

	It("should normalize the api grant dimension without changing the definition", func() {
		permission := map[string]any{"target_app_code": "app1", "grant_dimension": "api"}
		mgr = newManager(map[string]any{"grant_permissions": []any{permission}})
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/permissions/grant/").
			AddMatcher(captureBody).
			Reply(201).
			JSON(map[string]any{"data": nil})

		Expect(mgr.GrantPermissions(ctx)).To(Succeed())
		Expect(bodies).To(Equal([]map[string]any{
			{"target_app_code": "app1", "grant_dimension": "gateway"},
		}))
		Expect(permission["grant_dimension"]).To(Equal("api"))
	})

	DescribeTable("should use bk_app_code only as a fallback when granting permissions",
		func(permission map[string]any, targetAppCode string) {
			original := maps.Clone(permission)
			mgr = newManager(map[string]any{"grant_permissions": []any{permission}})
			gock.New(endpoint).
				Post("/api/v2/sync/gateways/testing/permissions/grant/").
				AddMatcher(captureBody).
				Reply(201).
				JSON(map[string]any{"data": nil})

			Expect(mgr.GrantPermissions(ctx)).To(Succeed())
			Expect(bodies).To(Equal([]map[string]any{
				{"target_app_code": targetAppCode, "grant_dimension": "gateway"},
			}))
			Expect(permission).To(Equal(original))
		},
		Entry("only the old field is present",
			map[string]any{"bk_app_code": "old-app", "grant_dimension": "api"}, "old-app"),
		Entry(
			"the new field takes precedence",
			map[string]any{
				"bk_app_code":     "old-app",
				"target_app_code": "new-app",
				"grant_dimension": "api",
			},
			"new-app",
		),
	)

	It("should apply permissions to the defined gateway with defaults", func() {
		gock.New(endpoint).
			Post("/api/v2/open/gateways/other-gateway/permissions/apply/").
			AddMatcher(captureBody).
			Reply(200).
			JSON(map[string]any{"data": map[string]any{"record_id": 1}})

		result, err := mgr.ApplyPermissions(ctx)
		Expect(err).To(BeNil())
		Expect(result["result_0"]).To(Equal(map[string]any{"record_id": float64(1)}))
		Expect(bodies).To(Equal([]map[string]any{
			{"target_app_code": "my-app", "applicant": "my-app", "grant_dimension": "gateway"},
		}))
	})

	It("should add related apps", func() {
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/related-apps/").
			AddMatcher(captureBody).
			Reply(201).
			JSON(map[string]any{"data": nil})

		Expect(mgr.AddRelatedApps(ctx)).To(Succeed())
		Expect(bodies).To(Equal([]map[string]any{
			{"related_app_codes": []any{"app1"}},
		}))
	})

	It("should skip adding related apps when there are none", func() {
		// gock fails the request if one is sent
		Expect(newManager(map[string]any{"related_apps": nil}).AddRelatedApps(ctx)).To(Succeed())
	})

	DescribeTable("should accept related apps from an in-memory definition",
		func(apps any) {
			mgr = newManager(map[string]any{"related_apps": apps})
			gock.New(endpoint).Post("/api/v2/sync/gateways/testing/related-apps/").
				AddMatcher(captureBody).Reply(201).BodyString(`{"data":null}`)
			Expect(mgr.AddRelatedApps(ctx)).To(Succeed())
			Expect(bodies).To(Equal([]map[string]any{
				{"related_app_codes": []any{"app1", "app2"}},
			}))
		},
		Entry("string slice", []string{"app1", "app2"}),
		Entry("YAML slice", []any{"app1", "app2"}),
	)

	DescribeTable("should reject invalid related apps before sending a request",
		func(apps any, message string) {
			mgr = newManager(map[string]any{"related_apps": apps})
			Expect(mgr.AddRelatedApps(ctx)).To(MatchError(ContainSubstring(message)))
		},
		Entry("string", "app1", "related_apps must be a list of app codes"),
		Entry("object", map[string]any{"app": "app1"}, "related_apps must be a list of app codes"),
		Entry("non-string element", []any{"app1", 2}, "related_apps must be a list of app codes"),
	)

	DescribeTable("should skip empty related app lists",
		func(apps any) {
			Expect(newManager(map[string]any{"related_apps": apps}).AddRelatedApps(ctx)).To(Succeed())
		},
		Entry("string slice", []string{}),
		Entry("YAML slice", []any{}),
	)

	It("should sync the basic info, stages and resources", func() {
		mgr = newManager(map[string]any{
			"apigateway": map[string]any{"description": "my gateway"},
			"stages":     []any{map[string]any{"name": "prod"}},
		})
		gock.New(endpoint).Post("/api/v2/sync/gateways/testing/").AddMatcher(captureBody).
			Reply(200).JSON(map[string]any{"data": map[string]any{"id": 1}})
		gock.New(endpoint).Post("/api/v2/sync/gateways/testing/stages/").AddMatcher(captureBody).
			Reply(200).JSON(map[string]any{"data": map[string]any{"id": 2}})
		gock.New(endpoint).Post("/api/v2/sync/gateways/testing/resources/").AddMatcher(captureBody).
			Reply(200).JSON(map[string]any{"data": map[string]any{"added": []any{}}})

		info, err := mgr.SyncBasicInfo(ctx)
		Expect(err).To(BeNil())
		Expect(info).To(Equal(map[string]any{"id": float64(1)}))

		stages, err := mgr.SyncStagesConfig(ctx)
		Expect(err).To(BeNil())
		Expect(stages).To(Equal(map[string]any{"prod": map[string]any{"id": float64(2)}}))

		_, err = mgr.SyncResourcesConfig(ctx, map[string]any{"content": "paths: {}", "delete": true})
		Expect(err).To(BeNil())

		Expect(bodies).To(Equal([]map[string]any{
			{"description": "my gateway"},
			{"name": "prod"},
			{"content": "paths: {}", "delete": true},
		}))
	})

	It("should upload the resource docs", func() {
		dir := GinkgoT().TempDir()
		Expect(os.WriteFile(filepath.Join(dir, "doc.md"), []byte("# doc"), 0o600)).To(Succeed())
		mgr = newManager(map[string]any{
			"resource_docs": map[string]any{"basedir": dir},
		})

		var contentType string
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/resource-docs/").
			AddMatcher(func(req *http.Request, _ *gock.Request) (bool, error) {
				contentType = req.Header.Get("Content-Type")
				return true, nil
			}).
			Reply(200).
			JSON(map[string]any{"data": nil})

		Expect(mgr.SyncResourceDocByArchive(ctx)).To(Succeed())
		Expect(contentType).To(HavePrefix("multipart/form-data"))
	})

	It("should get the public key", func() {
		gock.New(endpoint).
			Get("/api/v2/sync/gateways/testing/public_key/").
			Reply(200).
			JSON(map[string]any{"data": map[string]any{"public_key": "my-key"}})

		publicKey, err := mgr.GetPublicKeyString(ctx)
		Expect(err).To(BeNil())
		Expect(publicKey).To(Equal("my-key"))
	})

	It("should return the error of a response that is not json", func() {
		gock.New(endpoint).
			Get("/api/v2/sync/gateways/testing/resource_versions/latest/").
			Reply(200).
			BodyString("<html>login</html>")

		_, err := mgr.GetLatestResourceVersion(ctx)
		Expect(err).To(MatchError(ContainSubstring("decode response")))
	})

	It("should release with the comment", func() {
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/resource_versions/release/").
			AddMatcher(captureBody).
			Reply(200).
			JSON(map[string]any{
				"data": map[string]any{"version": "1.0.0+prod", "stage_names": []any{"prod"}},
			})

		_, err := mgr.Release(ctx, "1.0.0+prod", "release comment")
		Expect(err).To(BeNil())
		Expect(bodies).To(Equal([]map[string]any{
			{"version": "1.0.0+prod", "stage_names": []any{"prod"}, "comment": "release comment"},
		}))
	})
})
