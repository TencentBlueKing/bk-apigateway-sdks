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

package util

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/go-openapi/spec"
)

func MergeSwaggerConfig(swagger spec.Swagger, routeMap map[string]*RouteConfig) spec.Swagger {
	for path, pathItem := range swagger.Paths.Paths {
		operationMap := GetPathOperationMap(pathItem)
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
		data, _ = JsonToYAML(data)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetMcpToolAndValidate 获取mcp工具并校验
func GetMcpToolAndValidate(swagger spec.Swagger, routeMap map[string]*RouteConfig, tools []string) []string {
	// 是否有指定tool
	isSpecified := len(tools) > 0
	toolMap := make(map[string]struct{})       // 指定的tool的资源
	canUseToolMap := make(map[string]struct{}) // 可以作为mcp tool 的资源
	allToolMap := make(map[string]struct{})    // 所有资源
	for _, tool := range tools {
		routeConfig, ok := routeMap[tool]
		if !ok {
			log.Fatalf("tool: %s not found in route configs", tool)
			return []string{}
		}
		// 判断是否有开启mcp
		if !routeConfig.Config.EnableMcp {
			log.Fatalf("tool: %s not enable mcp", tool)
			return []string{}
		}
		toolMap[tool] = struct{}{}
	}

	// Iterate through all paths in the Swagger specification
	for path, pathItem := range swagger.Paths.Paths {
		operationMap := GetPathOperationMap(pathItem)
		// Process each operation in the operation map
		for method, operation := range operationMap {
			// 构造匹配键
			key := fmt.Sprintf("%s:%s", path, method)
			// 如果配置了mcp,则校验
			if c, exists := routeMap[key]; exists && c.Config != nil && c.Config.EnableMcp {
				// 如果指定了operationID,则使用指定的
				if operation.ID == "" {
					operation.ID = c.OperationID
				}
				allToolMap[operation.ID] = struct{}{}
				_, ok := toolMap[operation.ID]
				// 如果指定了tool,不在指定的tool中,则跳过校验
				if isSpecified && !ok {
					continue
				}

				if len(operation.Parameters) > 0 {
					// 如果没有指定tool,则可以作为mcp tool
					if !isSpecified {
						tools = append(tools, operation.ID)
					} else {
						canUseToolMap[operation.ID] = struct{}{}
					}
					continue
				}
				// 如果没有参数,开启mcp,则校验 nonSchema 是否设置
				if !c.Config.NonSchema {
					log.Fatalf("tool: %s enables mcp, but nonSchema is not set", operation.ID)
					return []string{}
				}
				// 如果没有指定tool,则可以作为mcp  tool
				if !isSpecified {
					tools = append(tools, operation.ID)
				} else {
					canUseToolMap[operation.ID] = struct{}{}
				}
			}
		}
	}
	if isSpecified {
		// 如果指定了tool,则校验是否都在canUseToolMap中
		for tool := range toolMap {
			// 判断是否在所有的tool中
			_, ok := allToolMap[tool]
			if !ok {
				log.Fatalf("tool: %s not found in all resources", tool)
				return []string{}
			}
			// 判断是否在canUseToolMap,不在则报错
			_, ok = canUseToolMap[tool]
			if !ok {
				log.Fatalf("tool: %s not found in all resources", tool)
				return []string{}
			}
			// 判断是否在canUseToolMap,不在则报错
			_, ok = canUseToolMap[tool]
			if !ok {
				log.Fatalf("tool: %s not found in can used resources, please check nonSchema or mcp enabled", tool)
				return []string{}
			}
		}
	}
	return tools
}

// GetPathOperationMap 获取path的operation map
func GetPathOperationMap(pathItem spec.PathItem) map[string]*spec.Operation {
	operationMap := make(map[string]*spec.Operation)
	// Populate the operation map with all available HTTP methods
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
	return operationMap
}
