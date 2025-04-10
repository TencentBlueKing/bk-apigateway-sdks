package router

import (
	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_fram/example/api"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_fram/middleware"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_fram/model"
)

func New() *gin.Engine {
	r := gin.Default()
	r.GET("/testapi/get-string-by-int/:some_id", middleware.WithGatewayResourceConfig(
		model.APIGatewayResourceConfig{
			IsPublic:             true,
			AllowApplyPermission: true,
			MatchSubpath:         false,
			EnableWebsocket:      false,
			Backend: model.BackendConfig{
				Path:   "/testapi/get-string-by-int/{some_id}",
				Method: "get",
			},
			PluginConfigs: []*model.PluginConfig{
				model.BuildPluginConfigWithType(model.PluginTypeBKCors, model.CorsConfig{
					AllowOrigins:    "*",
					AllowMethods:    "**",
					AllowHeaders:    "**",
					ExposeHeaders:   "",
					MaxAge:          0,
					AllowCredential: false,
				}),
			},
			AuthConfig: model.AuthConfig{},
		},
	), api.GetStringByInt)
	return r
}
