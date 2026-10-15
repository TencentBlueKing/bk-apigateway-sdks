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
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/spec"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/util"
)

// GenResourceYamlFromSwaggerJson 生成资源配置 yaml
// 从 swagger.json 文件生成资源配置 yaml
func GenResourceYamlFromSwaggerJson(docPath string, engine *gin.Engine, subPath string) string {
	// 获取 route 网关配置
	routeConfigMap := util.GetRouteConfigMap(engine)
	// 解析 Swagger 文件
	data, err := os.ReadFile(docPath)
	if err != nil {
		log.Fatal(err)
	}
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
