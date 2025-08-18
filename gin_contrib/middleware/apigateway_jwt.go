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
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/util"
	"github.com/TencentBlueKing/bk-apigateway-sdks/manager"
)

var (
	once              sync.Once
	publicMemoryCache *manager.PublicKeyMemoryCache
)

const (
	BkGatewayJWTHeaderKey = "X-Bkapi-Jwt"
)

func init() {
	once.Do(func() {
		config := bkapi.ClientConfig{}
		publicMemoryCache = manager.NewDefaultPublicKeyMemoryCache(config)
	})
}

// GatewayJWTAuthMiddleware 网关JWT鉴权中间件: 用于校验网关JWT,会进行应用、用户认证结构校验，需要使用网关
// RegisterBkAPIGatewayRoute 或者 RegisterBkAPIGatewayRouteWithGroup 注册路由
func GatewayJWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		signedToken := c.GetHeader(BkGatewayJWTHeaderKey)
		if signedToken == "" {
			util.UnauthorizedJSONResponse(c, "no authorization credentials provided")
			c.Abort()
			return
		}
		// get public key
		parser := manager.NewRsaJwtTokenParser(publicMemoryCache)
		claims, err := parser.Parse(signedToken)
		if err != nil {
			util.UnauthorizedJSONResponse(c, "token is invalid")
			c.Abort()
			return
		}
		// 获取route网关配置:默认校验应用是否通过认证
		config := util.GetRouteConfig(c.FullPath(), c.Request.Method)
		if (config != nil && config.AuthConfig.AppVerifiedRequired && !claims.App.Verified) || config == nil {
			util.UnauthorizedJSONResponse(c, fmt.Sprintf("app: %s is not verified", claims.App.BkAppCode))
			c.Abort()
			return
		}
		if config != nil && config.AuthConfig.UserVerifiedRequired && !claims.User.Verified {
			util.UnauthorizedJSONResponse(c, fmt.Sprintf("user: %s is not verified", claims.User.Username))
			c.Abort()
			return
		}

		if claims.App != nil && claims.App.BkAppCode != "" {
			util.SetJwtAppCode(c, claims.App.BkAppCode)
		}

		if claims.App != nil && claims.App.AppCode != "" {
			util.SetJwtAppCode(c, claims.App.AppCode)
		}

		if claims.User != nil && claims.User.Username != "" {
			util.SetJwtUserName(c, claims.User.Username)
		}

		c.Next()
	}
}
