package router

import (
	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/example/api"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/middleware"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/model"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/util"
)

func New() *gin.Engine {
	r := gin.Default()
	basicConfig := model.ResourceBasicConfig{
		IsPublic:             true,
		AllowApplyPermission: true,
		MatchSubpath:         false,
		EnableWebsocket:      false,
	}
	// 共用的插件配置
	headerWriterPlugin := model.BuildResourcePluginConfigWithType(
		model.PluginTypeHeaderRewrite, model.HeaderRewriteConfig{
			Set:    []model.HeaderRewriteValue{{Key: "X-Test", Value: "test"}},
			Remove: []model.HeaderRewriteValue{{Key: "X-Test2"}},
		})
	util.RegisterBkAPIGatewayRoute(r, "POST", "/testapi/update-product/:product_id",
		model.NewAPIGatewayResourceConfig(
			basicConfig,
			basicConfig.WithBackend(model.BackendConfig{
				Path:   "/testapi/update-product/{product_id}",
				Method: "post",
			}),
			basicConfig.WithPluginConfig(
				headerWriterPlugin,
			)),
		api.UpdateProduct)

	// group
	petGroup := r.Group("/testapi/pets")
	util.RegisterBkAPIGatewayRouteWithGroup(petGroup, "GET", "/:id/",
		model.NewAPIGatewayResourceConfig(
			basicConfig,
			basicConfig.WithBackend(model.BackendConfig{
				Path:   "/testapi/pets/{id}/",
				Method: "get",
			}),
			// 覆盖基础的通用配置
			basicConfig.WithPublic(false),
			basicConfig.WithPluginConfig(
				headerWriterPlugin,
				// 独自的插件配置
				model.BuildResourcePluginConfigWithType(model.PluginTypeBKCors, model.CorsConfig{
					AllowOrigins:    "*",
					AllowMethods:    "**",
					AllowHeaders:    "**",
					ExposeHeaders:   "",
					MaxAge:          0,
					AllowCredential: false,
				}))),
		// 网关jwt中间件，校验用户和应用是否合法
		middleware.GatewayJWTAuthMiddleware(),
		api.GetPetByID)
	return r
}
