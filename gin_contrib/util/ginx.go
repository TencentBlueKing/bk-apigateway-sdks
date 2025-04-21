package util

import (
	"fmt"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/model"
)

const (
	BkGatewayJwtClaimsUserName = "X-Bkapi-Jwt-Username"
	BkGatewayJwtClaimsAppCode  = "X-Bkapi-Jwt-Appcode"
)

// apiResourceConfigs 用于网关相关配置
var apiResourceConfigs = make(map[string]*model.APIGatewayResourceConfig)

func RegisterBkAPIGatewayRoute(
		engine *gin.Engine,
		method string,
		path string,
		config model.APIGatewayResourceConfig,
		handlers ...gin.HandlerFunc,
) {
	// 生成标准化的路由标识
	normalizedPath := ConvertExpressPathToSwagger(path)
	key := fmt.Sprintf("%s:%s", normalizedPath, strings.ToLower(method))
	apiResourceConfigs[key] = &config
	engine.Handle(method, path, handlers...)
}

func RegisterBkAPIGatewayRouteWithGroup(
		group *gin.RouterGroup,
		method string,
		path string,
		config model.APIGatewayResourceConfig,
		handlers ...gin.HandlerFunc,
) {
	// 生成标准化的路由标识
	normalizedPath := ConvertExpressPathToSwagger(joinPaths(group.BasePath(), path))
	key := fmt.Sprintf("%s:%s", normalizedPath, strings.ToLower(method))
	apiResourceConfigs[key] = &config
	group.Handle(method, path, handlers...)
}

func GetRouteConfigMap(engine *gin.Engine) map[string]*model.APIGatewayResourceConfig {
	result := make(map[string]*model.APIGatewayResourceConfig)
	for _, route := range engine.Routes() {
		// 生成标准路由标识
		normalizedPath := ConvertExpressPathToSwagger(route.Path)
		key := fmt.Sprintf("%s:%s", normalizedPath, strings.ToLower(route.Method))
		// 直接从同步Map获取配置
		if val, ok := apiResourceConfigs[key]; ok {
			result[key] = val
		}
	}
	return result
}

func GetRouteConfig(path string, method string) *model.APIGatewayResourceConfig {
	normalizedPath := ConvertExpressPathToSwagger(path)
	return apiResourceConfigs[fmt.Sprintf("%s:%s", normalizedPath, strings.ToLower(method))]
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

// SetJwtUserName 设置jwt user信息到context
func SetJwtUserName(c *gin.Context, username string) {
	c.Set(BkGatewayJwtClaimsUserName, username)
}

// SetJwtAppCode 获取jwt appCode信息到context
func SetJwtAppCode(c *gin.Context, appCode string) {
	c.Set(BkGatewayJwtClaimsAppCode, appCode)
}

// GetJwtAppCode GetJwtInfo 获取jwt信息
func GetJwtAppCode(c *gin.Context) string {
	return c.GetString(BkGatewayJwtClaimsAppCode)
}

// GetJwtUserName 获取username信息
func GetJwtUserName(c *gin.Context) string {
	return c.GetString(BkGatewayJwtClaimsUserName)
}
