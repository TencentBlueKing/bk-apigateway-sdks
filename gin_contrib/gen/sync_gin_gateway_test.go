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
		Stage: model.StageConfig{
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
	definitionConfig := GenDefinitionYaml(config)
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
