package middleware

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/util"
	"github.com/TencentBlueKing/bk-apigateway-sdks/manager"
)

type JwtConfig struct {
	CheckUser bool // 是否校验用户 verified
	CheckApp  bool // 是否校验app verified
}

var once sync.Once
var publicMemoryCache *manager.PublicKeyMemoryCache

const (
	BkGatewayJWTHeaderKey = "X-Bkapi-Jwt"
	BkGatewayJwtClaims    = "X-Bkapi-Jwt-Claims"
)

func init() {
	once.Do(func() {
		config := bkapi.ClientConfig{}
		publicMemoryCache = manager.NewDefaultPublicKeyMemoryCache(config)
	})

}

func GatewayJWTAuthMiddleware(config JwtConfig) func(c *gin.Context) {
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
		if config.CheckUser && !claims.User.Verified {
			util.UnauthorizedJSONResponse(c, fmt.Sprintf("user:%s is not verified", claims.User.Username))
			c.Abort()
			return
		}
		c.Set(BkGatewayJwtClaims, claims)
		c.Next()
	}
}
