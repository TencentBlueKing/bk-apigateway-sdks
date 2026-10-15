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

package gen

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/example/router"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/model"
)

func TestSyncGinGateway(t *testing.T) {
	// 初始化配置
	config := &model.APIConfig{
		Release: model.ReleaseConfig{
			Version: "1.0.0+prod",
			Title:   "初始版本",
			Comment: "首次发布",
		},
		APIGateway: model.GatewayConfig{
			Description:   "示例网关",
			DescriptionEn: "Example Gateway",
			IsPublic:      true,
			IsOfficial:    false,
			Maintainers:   []string{"handryhan"},
		},
		Stage: &model.StageConfig{
			Name:           "prod",
			Description:    "生产环境",
			DescriptionEn:  "Production",
			BackendSubPath: "/api",
			BackendTimeout: 30,
			BackendHost:    "http://api.example.com",
			PluginConfigs: []*model.PluginConfig{
				model.BuildStagePluginConfigWithType(
					model.PluginTypeHeaderRewrite,
					model.HeaderRewriteConfig{
						Set: []model.HeaderRewriteValue{
							{Key: "X-Real-IP", Value: "test"},
						},
						Remove: []model.HeaderRewriteValue{
							{Key: "X-Forwarded-For"},
						},
					}),
			},
		},
		GrantPermissions: model.GrantPermissionConfig{
			GatewayApps: []string{"app1"},
			ResourceApps: map[string][]string{
				"app2": {"get_pet_by_id"},
			},
		},
		RelatedApps: []string{"myapp"},
		ResourceDocs: model.ResourceDocConfig{
			BaseDir: "../example/docs/",
		},
	}

	// 生成定义配置
	definitionConfig := GenDefinitionYaml(config, "../example/docs/swagger.json", router.New())
	definitionFilePath := filepath.Join("./example", "definition.yaml")
	// 先创建目录（递归创建）
	if err := os.MkdirAll(filepath.Dir(definitionFilePath), 0o755); err != nil {
		t.Fatal("创建目录失败: " + err.Error())
	}
	err := os.WriteFile(definitionFilePath, []byte(definitionConfig), 0o644)
	if err != nil {
		t.Fatal(err)
	}
	// 生成 resource 配置
	resourcesFilePath := filepath.Join("./example", "resources.yaml")
	resourcesYaml := GenResourceYamlFromSwaggerJson("../example/docs/swagger.json", router.New(), "")
	if err := os.WriteFile(resourcesFilePath, []byte(resourcesYaml), 0o644); err != nil {
		t.Fatal(err)
	}
}

type recordedRequest struct {
	Method   string
	Path     string
	Query    string
	Body     map[string]any
	ZipFiles []string
}

type mockV2Gateway struct {
	mu            sync.Mutex
	requests      []recordedRequest
	versionExists bool
}

func (g *mockV2Gateway) reply(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": data})
}

func (g *mockV2Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	record := recordedRequest{Method: r.Method, Path: r.URL.Path, Query: r.URL.RawQuery}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		file, _, err := r.FormFile("file")
		if err == nil {
			content, _ := io.ReadAll(file)
			if reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content))); err == nil {
				for _, f := range reader.File {
					record.ZipFiles = append(record.ZipFiles, f.Name)
				}
			}
		}
	} else if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&record.Body)
	}
	g.mu.Lock()
	g.requests = append(g.requests, record)
	g.mu.Unlock()

	prefix := "/api/bk-apigateway/prod/api/v2/sync/gateways/testing/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{"code": "NOT_FOUND", "message": r.URL.Path},
		})
		return
	}

	switch sub := strings.TrimPrefix(r.URL.Path, prefix); {
	case sub == "":
		g.reply(w, http.StatusOK, map[string]any{"id": 1, "name": "testing"})
	case sub == "stages/":
		g.reply(w, http.StatusOK, map[string]any{"id": 1, "name": "prod"})
	case sub == "resources/":
		g.reply(w, http.StatusOK, map[string]any{
			"added": []any{}, "updated": []any{}, "deleted": []any{},
		})
	case sub == "permissions/grant/", sub == "resource-docs/", sub == "related-apps/":
		g.reply(w, http.StatusCreated, nil)
	case sub == "resource_versions/" && r.Method == http.MethodGet:
		count := 0
		if g.versionExists {
			count = 1
		}
		g.reply(w, http.StatusOK, map[string]any{"count": count, "results": []any{}})
	case sub == "resource_versions/":
		g.reply(w, http.StatusOK, map[string]any{"id": 1, "version": record.Body["version"]})
	case sub == "resource_versions/release/":
		g.reply(w, http.StatusOK, map[string]any{
			"version": record.Body["version"], "stage_names": record.Body["stage_names"],
		})
	case sub == "stages/prod/mcp-servers/":
		g.reply(w, http.StatusOK, []any{map[string]any{"id": 1, "name": "mcp", "action": "created"}})
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (g *mockV2Gateway) find(t *testing.T, method, pathSuffix string) []recordedRequest {
	t.Helper()
	var result []recordedRequest
	for _, req := range g.requests {
		if req.Method == method && strings.HasSuffix(req.Path, pathSuffix) {
			result = append(result, req)
		}
	}
	if len(result) == 0 {
		t.Fatalf("no %s request to %s, requests: %+v", method, pathSuffix, g.requests)
	}
	return result
}

func TestSyncGinGatewayWithV2Apis(t *testing.T) {
	cases := []struct {
		name          string
		versionExists bool
		wantVersion   *regexp.Regexp
	}{
		{name: "new version", versionExists: false, wantVersion: regexp.MustCompile(`^1\.0\.0\+prod$`)},
		{name: "existing version", versionExists: true, wantVersion: regexp.MustCompile(`^1\.0\.0\+\d{14}$`)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gateway := &mockV2Gateway{versionExists: c.versionExists}
			baseDir, config := setupSyncGinGateway(t, gateway)

			SyncGinGateway(baseDir, "testing", config, true)

			resources := gateway.find(t, http.MethodPost, "/resources/")[0].Body
			if resources["doc_language"] != "zh" || resources["language"] != nil || resources["delete"] != true {
				t.Errorf("unexpected resources sync body: %v", resources)
			}

			// 和 apigw-manager 一样，同步流程不添加关联应用
			for _, req := range gateway.requests {
				if strings.HasSuffix(req.Path, "/related-apps/") {
					t.Errorf("unexpected request to add related apps: %+v", req)
				}
			}

			grant := gateway.find(t, http.MethodPost, "/permissions/grant/")[0].Body
			if grant["target_app_code"] != "app1" || grant["grant_dimension"] != "gateway" {
				t.Errorf("unexpected grant body: %v", grant)
			}

			docs := gateway.find(t, http.MethodPost, "/resource-docs/")[0]
			if len(docs.ZipFiles) != 1 || docs.ZipFiles[0] != "zh/get_pet.md" {
				t.Errorf("unexpected resource docs archive: %v", docs.ZipFiles)
			}

			lookup := gateway.find(t, http.MethodGet, "/resource_versions/")[0]
			if lookup.Query != "version=1.0.0%2Bprod" {
				t.Errorf("unexpected resource version lookup query: %s", lookup.Query)
			}

			created := gateway.find(t, http.MethodPost, "/resource_versions/")[0].Body
			version, _ := created["version"].(string)
			if !c.wantVersion.MatchString(version) || created["comment"] != "release comment" {
				t.Errorf("unexpected create resource version body: %v", created)
			}

			release := gateway.find(t, http.MethodPost, "/resource_versions/release/")[0].Body
			if release["version"] != version || release["comment"] != "release comment" {
				t.Errorf("unexpected release body: %v", release)
			}

			mcp := gateway.find(t, http.MethodPost, "/stages/prod/mcp-servers/")[0].Body
			servers, _ := mcp["mcp_servers"].([]any)
			if len(servers) != 1 {
				t.Fatalf("unexpected mcp servers sync body: %v", mcp)
			}
			want := map[string]any{
				"name":                           "mcp",
				"resource_names":                 []any{"get_pet"},
				"tool_names":                     []any{"fetch_pet"},
				"category_names":                 []any{"Official"},
				"oauth2_personal_client_enabled": true,
			}
			if !reflect.DeepEqual(servers[0], want) {
				t.Errorf("mcp server = %v, want %v", servers[0], want)
			}
		})
	}
}

// setupSyncGinGateway starts the gateway, and writes the definition and resources to the returned dir.
func setupSyncGinGateway(t *testing.T, gateway *mockV2Gateway) (string, *model.APIConfig) {
	t.Helper()
	server := httptest.NewServer(gateway)
	t.Cleanup(server.Close)
	t.Setenv("BK_API_URL_TMPL", server.URL+"/api/{api_name}")
	t.Setenv("BK_APP_CODE", "my-app")
	t.Setenv("BK_APP_SECRET", "secret")

	baseDir := t.TempDir()
	docsDir := filepath.Join(baseDir, "docs")
	if err := os.MkdirAll(filepath.Join(docsDir, "zh"), 0o755); err != nil {
		t.Fatal(err)
	}
	docFile := filepath.Join(docsDir, "zh", "get_pet.md")
	if err := os.WriteFile(docFile, []byte("# get pet"), 0o644); err != nil {
		t.Fatal(err)
	}
	definition := strings.Join([]string{
		"spec_version: 2",
		"apigateway:",
		"  description: testing",
		"stages:",
		"  - name: prod",
		"    mcp_servers:",
		"      - name: mcp",
		"        resource_names: [get_pet]",
		"        tool_names: [fetch_pet]",
		"        category_names: [Official]",
		"        oauth2_personal_client_enabled: true",
		"related_apps:",
		"  - app1",
		"  - app2",
		"grant_permissions:",
		"  - target_app_code: app1",
		"    grant_dimension: gateway",
		"resource_docs:",
		"  basedir: " + docsDir,
	}, "\n")
	if err := os.WriteFile(filepath.Join(baseDir, "definition.yaml"), []byte(definition), 0o644); err != nil {
		t.Fatal(err)
	}
	resourcesFile := filepath.Join(baseDir, "resources.yaml")
	if err := os.WriteFile(resourcesFile, []byte("swagger: '2.0'"), 0o644); err != nil {
		t.Fatal(err)
	}

	return baseDir, &model.APIConfig{
		Release:      model.ReleaseConfig{Version: "1.0.0+prod", Comment: "release comment"},
		Stage:        &model.StageConfig{Name: "prod", EnableMcpServers: true},
		ResourceDocs: model.ResourceDocConfig{BaseDir: docsDir, Language: "zh"},
	}
}
