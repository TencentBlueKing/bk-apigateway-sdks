package gen

import (
	"testing"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_fram/model"
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
		Stage: model.StageConfig{
			Name:           "prod",
			Description:    "生产环境",
			DescriptionEn:  "Production",
			BackendSubPath: "/api",
			BackendTimeout: 30,
			BackendHost:    "api.example.com",
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
	definitionConfig := GenDefinitionYaml(config)
	t.Log(definitionConfig)
}
