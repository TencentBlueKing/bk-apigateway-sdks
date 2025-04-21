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
	util.RegisterBkAPIGatewayRoute(r, "POST", "/testapi/update-product/:product_id",
		model.APIGatewayResourceConfig{
			IsPublic:             true,
			AllowApplyPermission: true,
			MatchSubpath:         false,
			EnableWebsocket:      false,
			Backend: model.BackendConfig{
				Path:   "/testapi/product/",
				Method: "POST",
			},
			PluginConfigs: []*model.PluginConfig{
				model.BuildResourcePluginConfigWithType(model.PluginTypeHeaderRewrite, model.HeaderRewriteConfig{
					Set:    []model.HeaderRewriteValue{{Key: "X-Test", Value: "test"}},
					Remove: []model.HeaderRewriteValue{{Key: "X-Test2"}},
				}),
			},
			AuthConfig: model.AuthConfig{},
		},
		api.UpdateProduct)

	// group
	petGroup := r.Group("/testapi/pets")
	util.RegisterBkAPIGatewayRouteWithGroup(petGroup, "GET", "/:id/",
		model.APIGatewayResourceConfig{
			IsPublic:             false,
			AllowApplyPermission: true,
			MatchSubpath:         false,
			EnableWebsocket:      false,
			Backend: model.BackendConfig{
				Path:   "/testapi/pets/{id}/",
				Method: "get",
			},
			PluginConfigs: []*model.PluginConfig{
				model.BuildResourcePluginConfigWithType(model.PluginTypeBKCors, model.CorsConfig{
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
		// 网关jwt中间件，校验用户和应用是否合法
		middleware.GatewayJWTAuthMiddleware(),
		api.GetPetByID)
	return r
}
