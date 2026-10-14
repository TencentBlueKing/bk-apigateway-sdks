package gen

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/spec"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/util"
)

// GenResourceYamlFromSwaggerJson 生成资源配置yaml
// 从 swagger.json文件生成资源配置yaml
func GenResourceYamlFromSwaggerJson(docPath string, engine *gin.Engine, subPath string) string {
	// 获取route 网关配置
	routeConfigMap := util.GetRouteConfigMap(engine)
	// 解析 Swagger 文件
	data, _ := os.ReadFile(docPath)
	var swagger spec.Swagger
	if err := json.Unmarshal(data, &swagger); err != nil {
		log.Fatal(err)
	}
	// 合并配置
	swagger = util.MergeSwaggerConfig(swagger, routeConfigMap, subPath)
	config, err := util.OutputResourceConfig(&swagger, "yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(config))
	return string(config)
}
