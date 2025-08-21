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
	"testing"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/example/router"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/model"
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
			APIType:       "public",
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
		RelatedApps: []string{"myapp"},
		ResourceDocs: model.ResourceDocConfig{
			BaseDir: "/data/docs",
		},
	}
	// 生成定义配置
	definitionConfig := GenDefinitionYaml(config, "../example/docs/swagger.json", router.New())
	t.Log(definitionConfig)
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
			APIType:       "public",
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
					Name:           "mcp-server-1",
					Description:    "mcp-server-1",
					IsPublic:       false,
					Status:         1,
					Labels:         []string{"label1", "label2"},
					TargetAppCodes: []string{"app1", "app2"},
					Tools:          []string{"update_product_set"},
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
	t.Log(definitionConfig)
}
