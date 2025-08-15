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
	"os"
	"path/filepath"
	"testing"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/example/router"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/model"
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
			APIType:       "10",
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
	// 生成resource配置
	resourcesFilePath := filepath.Join("./example", "resources.yaml")
	resourcesYaml := GenResourceYamlFromSwaggerJson("../example/docs/swagger.json", router.New())
	err = os.WriteFile(resourcesFilePath, []byte(resourcesYaml), 0o644)
	// SyncGinGateway(
	//	"./example/",
	//	"custom-gateway-go-demo", config, true)
}
