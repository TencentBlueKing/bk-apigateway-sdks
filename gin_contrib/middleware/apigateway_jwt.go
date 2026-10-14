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

package middleware

import (
	"cmp"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/util"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/manager"
)

const (
	BkGatewayJWTHeaderKey = "X-Bkapi-Jwt"
)

// 所有中间件共用一份公钥缓存，第一次创建中间件时才读取环境变量
var publicKeyCache = sync.OnceValue(func() *manager.PublicKeyMemoryCache {
	return manager.NewDefaultPublicKeyMemoryCache(manager.ConfigFromEnv())
})

// GatewayJWTAuthMiddleware 校验请求头中的网关 JWT，并检查应用和用户是否通过认证。
// 路由需要用 RegisterBkAPIGatewayRoute 或 RegisterBkAPIGatewayRouteWithGroup 注册。
//
// 请在环境变量加载之后再调用，例如在 main 里读完 .env 之后。
func GatewayJWTAuthMiddleware() func(c *gin.Context) {
	parser := manager.NewRsaJwtTokenParser(publicKeyCache())
	return func(c *gin.Context) {
		signedToken := c.GetHeader(BkGatewayJWTHeaderKey)
		if signedToken == "" {
			util.UnauthorizedJSONResponse(c, "no authorization credentials provided")
			c.Abort()
			return
		}
		claims, err := parser.Parse(signedToken)
		if err != nil {
			util.UnauthorizedJSONResponse(c, "token is invalid")
			c.Abort()
			return
		}
		// 没有应用或用户时视为未认证，应用优先取 bk_app_code
		app := cmp.Or(claims.App, &manager.ApigatewayJwtApp{})
		user := cmp.Or(claims.User, &manager.ApigatewayJwtUser{})
		appCode := cmp.Or(app.BkAppCode, app.AppCode)
		// 获取route网关配置:默认校验应用是否通过认证
		config := util.GetRouteConfig(c.FullPath(), c.Request.Method)
		if (config != nil && config.AuthConfig.AppVerifiedRequired && !app.Verified) || config == nil {
			util.UnauthorizedJSONResponse(c, fmt.Sprintf("app: %s is not verified", appCode))
			c.Abort()
			return
		}
		if config.AuthConfig.UserVerifiedRequired && !user.Verified {
			util.UnauthorizedJSONResponse(c, fmt.Sprintf("user: %s is not verified", user.Username))
			c.Abort()
			return
		}

		if appCode != "" {
			util.SetJwtAppCode(c, appCode)
		}
		if user.Username != "" {
			util.SetJwtUserName(c, user.Username)
		}

		c.Next()
	}
}
