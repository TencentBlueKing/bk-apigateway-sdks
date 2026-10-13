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
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/example/router"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/model"
)

func TestGenDefinitionConfig(t *testing.T) {
	// 初始化配置
	config := &model.APIConfig{
		Release: model.ReleaseConfig{
			Version: "1.0.0",
			Title:   "初始版本",
			Comment: "首次发布",
		},
		APIGateway: model.GatewayConfig{
			Description:   "示例网关",
			DescriptionEn: "Example Gateway",
			IsPublic:      true,
			IsOfficial:    true,
			Maintainers:   []string{"user1", "user2"},
		},
		Stage: &model.StageConfig{
			Name:           "prod",
			Description:    "生产环境",
			DescriptionEn:  "Production",
			BackendSubPath: "/api",
			BackendTimeout: 30,
			BackendHost:    "api.example.com",
			PluginConfigs: []*model.PluginConfig{
				model.BuildStagePluginConfigWithType(
					model.PluginTypeHeaderRewrite,
					model.HeaderRewriteConfig{
						Set: []model.HeaderRewriteValue{
							{Key: "X-Real-IP", Value: "123"},
						},
						Remove: []model.HeaderRewriteValue{
							{Key: "X-Forwarded-For"},
						},
					}),
			},
			EnvVars: map[string]string{
				"foo": "bar",
			},
		},
		GrantPermissions: model.GrantPermissionConfig{
			GatewayApps: []string{"app1"},
			ResourceApps: map[string][]string{
				"app2": {"res1", "res2"},
			},
		},
		RelatedApps: []string{"myapp", "otherapp"},
		ResourceDocs: model.ResourceDocConfig{
			BaseDir: "/data/docs",
		},
	}
	// 生成定义配置
	definitionConfig := GenDefinitionYaml(config, "../example/docs/swagger.json", router.New())
	var parsed struct {
		Gateway     map[string]any   `yaml:"apigateway"`
		Permissions []map[string]any `yaml:"grant_permissions"`
		RelatedApps []string         `yaml:"related_apps"`
	}
	if err := yaml.Unmarshal([]byte(definitionConfig), &parsed); err != nil {
		t.Fatal(err)
	}
	expected := []map[string]any{
		{"target_app_code": "app1", "grant_dimension": "gateway"},
		{"target_app_code": "app2", "grant_dimension": "resource", "resource_names": []any{"res1", "res2"}},
	}
	if !reflect.DeepEqual(parsed.Permissions, expected) {
		t.Fatalf("invalid v2 permissions: %#v", parsed.Permissions)
	}
	if !reflect.DeepEqual(parsed.RelatedApps, []string{"myapp", "otherapp"}) {
		t.Errorf("invalid related apps: %v", parsed.RelatedApps)
	}
	wantGateway := map[string]any{
		"description":    "示例网关",
		"description_en": "Example Gateway",
		"is_public":      true,
		"api_type":       1,
		"maintainers":    []any{"user1", "user2"},
	}
	if !reflect.DeepEqual(parsed.Gateway, wantGateway) {
		t.Errorf("gateway = %v, want %v", parsed.Gateway, wantGateway)
	}
}

func TestGenDefinitionConfigWithMcpServer(t *testing.T) {
	// 初始化配置
	config := &model.APIConfig{
		Release: model.ReleaseConfig{
			Version: "1.0.0",
			Title:   "初始版本",
			Comment: "首次发布",
		},
		APIGateway: model.GatewayConfig{
			Description:   "示例网关",
			DescriptionEn: "Example Gateway",
			IsPublic:      true,
			IsOfficial:    false,
			Maintainers:   []string{"user1", "user2"},
		},
		Stage: &model.StageConfig{
			Name:           "prod",
			Description:    "生产环境",
			DescriptionEn:  "Production",
			BackendSubPath: "/api",
			BackendTimeout: 30,
			BackendHost:    "api.example.com",
			PluginConfigs: []*model.PluginConfig{
				model.BuildStagePluginConfigWithType(
					model.PluginTypeHeaderRewrite,
					model.HeaderRewriteConfig{
						Set: []model.HeaderRewriteValue{
							{Key: "X-Real-IP", Value: "123"},
						},
						Remove: []model.HeaderRewriteValue{
							{Key: "X-Forwarded-For"},
						},
					}),
			},
			EnvVars: map[string]string{
				"foo": "bar",
			},
			EnableMcpServers: true,
			McpServerConfigs: []*model.McpServer{
				{
					Name:                        "mcp-server-1",
					Title:                       "MCP服务1",
					Description:                 "mcp-server-1",
					IsPublic:                    false,
					ProtocolType:                model.MCPServerProtocolSSE,
					Status:                      1,
					Labels:                      []string{"label1", "label2"},
					TargetAppCodes:              []string{"app1", "app2"},
					ResourceNames:               []string{"update_product_set"},
					ToolNames:                   []string{"update_product"},
					Oauth2PublicClientEnabled:   true,
					Oauth2PersonalClientEnabled: true,
					RawResponseEnabled:          false,
					CategoryNames:               []string{"Official", "Automation"},
				},
			},
		},
		GrantPermissions: model.GrantPermissionConfig{
			GatewayApps: []string{"app1"},
			ResourceApps: map[string][]string{
				"app2": {"res1", "res2"},
			},
		},
		RelatedApps: []string{"myapp"},
		ResourceDocs: model.ResourceDocConfig{
			BaseDir: "/data/docs",
		},
	}
	// 生成定义配置
	definitionConfig := GenDefinitionYaml(config, "../example/docs/swagger.json", router.New())
	var parsed struct {
		Stages []struct {
			McpServers []map[string]any `yaml:"mcp_servers"`
		} `yaml:"stages"`
	}
	if err := yaml.Unmarshal([]byte(definitionConfig), &parsed); err != nil {
		t.Fatal(err)
	}
	if len(parsed.Stages) != 1 || len(parsed.Stages[0].McpServers) != 1 {
		t.Fatalf("invalid mcp servers: %s", definitionConfig)
	}
	server := parsed.Stages[0].McpServers[0]
	if server["oauth2_public_client_enabled"] != true || server["oauth2_personal_client_enabled"] != true {
		t.Errorf("invalid oauth2 flags of the mcp server: %v", server)
	}
	if !reflect.DeepEqual(server["tool_names"], []any{"update_product"}) ||
		!reflect.DeepEqual(server["category_names"], []any{"Official", "Automation"}) {
		t.Errorf("invalid tool or category names of the mcp server: %v", server)
	}
}
