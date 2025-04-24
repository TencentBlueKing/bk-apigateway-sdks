package util

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/model"
)

const (
	BkGatewayJwtClaimsUserName = "X-Bkapi-Jwt-Username"
	BkGatewayJwtClaimsAppCode  = "X-Bkapi-Jwt-Appcode"
)

type RouteConfig struct {
	OperationID string
	Config      *model.APIGatewayResourceConfig
}

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

func GetRouteConfigMap(engine *gin.Engine) map[string]*RouteConfig {
	result := make(map[string]*RouteConfig)
	for _, route := range engine.Routes() {
		// 生成标准路由标识
		normalizedPath := ConvertExpressPathToSwagger(route.Path)
		key := fmt.Sprintf("%s:%s", normalizedPath, strings.ToLower(route.Method))
		routeConfig := &RouteConfig{
			OperationID: GenerateOperationID(route.Handler),
		}
		// 直接从同步Map获取配置
		if val, ok := apiResourceConfigs[key]; ok {
			routeConfig.Config = val
		}
		result[key] = routeConfig
	}
	return result
}

// GenerateOperationID 从 handler 路径生成蛇形 OperationID
// 示例输入: "github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/example/api.UpdateProduct"
// 示例输出: "api_update_product"
func GenerateOperationID(operation string) string {
	pkgName, funcName := splitPackageAndFunc(operation)
	return fmt.Sprintf("%s_%s", strings.ToLower(strings.ReplaceAll(pkgName, ".", "_")), toSnakeCase(funcName))
}

// splitPackageAndFunc 分离包名和函数名（支持多级包名）
func splitPackageAndFunc(s string) (pkg, funcName string) {
	// 处理多级包名
	if idx := strings.LastIndex(s, "/"); idx != -1 {
		s = s[idx+1:]
	}
	if idx := strings.LastIndex(s, "."); idx != -1 {
		return s[:idx], s[idx+1:]
	}
	return "", s
}

// toSnakeCase大驼峰转蛇形命名
func toSnakeCase(str string) string {
	// Step 1: 处理连续大写后跟小写 (如 HTTPHandler → HTTP_Handler)
	snake := regexp.MustCompile("([A-Z]+)([A-Z][a-z])").ReplaceAllString(str, "${1}_${2}")

	// Step 2: 处理小写/数字后跟大写 (如 UpdateProduct → Update_Product)
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")

	// Step 4: 处理数字后跟字母 (如 JSON2XML → JSON_2_XML)
	snake = regexp.MustCompile("([0-9])([A-Za-z])").ReplaceAllString(snake, "${1}_${2}")

	// 最终处理
	snake = strings.ToLower(snake)
	snake = regexp.MustCompile("__+").ReplaceAllString(snake, "_") // 合并连续下划线
	return strings.Trim(snake, "_")
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
