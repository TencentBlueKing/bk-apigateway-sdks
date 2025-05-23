package gen

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/spec"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/util"
)

// GenResourceYamlFromSwaggerJson 生成资源配置yaml
// 从swagger.json文件生成资源配置yaml
func GenResourceYamlFromSwaggerJson(docPath string, engine *gin.Engine) string {
	// 获取route 网关配置
	routeConfigMap := util.GetRouteConfigMap(engine)
	// 解析 Swagger 文件
	data, _ := os.ReadFile(docPath)
	var swagger spec.Swagger
	if err := json.Unmarshal(data, &swagger); err != nil {
		log.Fatal(err)
	}
	// 合并配置
	swagger = mergeSwaggerConfig(swagger, routeConfigMap)

	config, err := OutputResourceConfig(&swagger, "yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(config))
	return string(config)
}

func mergeSwaggerConfig(swagger spec.Swagger, routeMap map[string]*util.RouteConfig) spec.Swagger {
	for path, pathItem := range swagger.Paths.Paths {
		operationMap := make(map[string]*spec.Operation)
		if pathItem.PathItemProps.Get != nil {
			operationMap["get"] = pathItem.Get
		}
		if pathItem.PathItemProps.Post != nil {
			operationMap["post"] = pathItem.Post
		}
		if pathItem.PathItemProps.Put != nil {
			operationMap["put"] = pathItem.Put
		}
		if pathItem.PathItemProps.Delete != nil {
			operationMap["delete"] = pathItem.Delete
		}
		if pathItem.PathItemProps.Options != nil {
			operationMap["options"] = pathItem.Options
		}
		if pathItem.PathItemProps.Head != nil {
			operationMap["head"] = pathItem.Head
		}
		if pathItem.PathItemProps.Patch != nil {
			operationMap["patch"] = pathItem.Patch
		}
		for method, operation := range operationMap {
			// 构造匹配键
			key := fmt.Sprintf("%s:%s", path, method)
			// 合并配置
			if c, exists := routeMap[key]; exists {
				if c.Config != nil {
					c.Config.Backend.Method = strings.ToLower(method)
					if c.Config.Backend.Path == "" {
						c.Config.Backend.Path = path
					}
					if c.Config.Backend.Method == "" {
						c.Config.Backend.Method = strings.ToLower(method)
					}
					if operation.Extensions == nil {
						operation.Extensions = spec.Extensions{}
					}
					// 保持一致
					c.Config.Backend.MatchSubpath = c.Config.MatchSubpath
					operation.Extensions.Add("x-bk-apigateway-resource", c.Config)
				}
				// 使用生成的OperationID作为路由ID
				if operation.ID == "" {
					operation.ID = c.OperationID
				}
			}
		}
	}
	return swagger
}

// OutputResourceConfig 输出配置
// 支持json和yaml格式
func OutputResourceConfig(doc *spec.Swagger, format string) ([]byte, error) {
	var data []byte
	var err error
	switch format {
	case "json":
		data, err = json.MarshalIndent(doc, "", "  ")
	case "yaml":
		data, err = json.MarshalIndent(doc, "", "  ")
		data, _ = util.JsonToYAML(data)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}
