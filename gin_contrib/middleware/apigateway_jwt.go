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
			util.SetJwtAppCode(c, claims.User.Username)
		}

		c.Next()
	}
}
