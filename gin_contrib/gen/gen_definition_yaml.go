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

package gen

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/model"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/util"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/spec"
)

//go:embed definition.tpl
var configTemplate []byte

func GenDefinitionYaml(config *model.APIConfig, docPath string, engine *gin.Engine) string {
	// 创建模板
	tmpl, err := template.New("config").Funcs(template.FuncMap{
		"indent": func(n int, s string) string {
			pad := strings.Repeat(" ", n)
			return pad + strings.ReplaceAll(s, "\n", "\n"+pad)
		},
	}).Parse(string(configTemplate))
	if err != nil {
		panic(fmt.Sprintf("gen definition yaml error: %v", err))
	}

	// 如果开起mcp，则需要根据 route 配置校验tool以及填充tool配置(如果没有指定)
	if config.Stage.EnableMcp {
		// 获取route 网关配置
		routeConfigMap := util.GetRouteConfigMap(engine)
		// 解析 Swagger 文件
		data, _ := os.ReadFile(docPath)
		var swagger spec.Swagger
		if err := json.Unmarshal(data, &swagger); err != nil {
			log.Fatal(err)
		}
		for _, mcpServer := range config.Stage.McpServerConfigs {
			mcpServer.Tools = util.GetMcpToolAndValidate(swagger, routeConfigMap, mcpServer.Tools)
		}

	}

	// 渲染输出
	var result strings.Builder
	if err := tmpl.Execute(&result, config); err != nil {
		panic(fmt.Sprintf("gen definition yaml error: %v", err))
	}
	fmt.Println("gen definition yaml:\n" + result.String())
	return result.String()
}
