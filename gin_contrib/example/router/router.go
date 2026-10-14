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

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/example/api"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/middleware"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/model"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/util"
)

// New creates the example router, whose routes are registered with the gateway resource configs.
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
			basicConfig.WithOperationID("update_product_set"),
			basicConfig.WithBackend(model.BackendConfig{
				Path:   "/testapi/update-product/{product_id}",
				Method: "post",
			}),
			// 开启 mcp 功能
			basicConfig.WithMcpEnable(true),
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
		// 网关 jwt 中间件，校验用户和应用是否合法
		middleware.GatewayJWTAuthMiddleware(),
		api.GetPetByID)
	return r
}
