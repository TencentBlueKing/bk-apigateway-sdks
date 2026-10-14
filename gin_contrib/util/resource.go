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
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/go-openapi/spec"
)

// MergeSwaggerConfig merges the route configs into the swagger, and sets the backend of each operation.
// If subPath is not empty, the backend paths of the default backend are prefixed with {env.api_sub_path}.
func MergeSwaggerConfig(swagger spec.Swagger, routeMap map[string]*RouteConfig, subPath string) spec.Swagger {
	if swagger.Paths == nil {
		return swagger
	}
	// 没有子路径时 definition.yaml 不定义 api_sub_path
	prefix := ""
	if strings.Trim(subPath, "/") != "" {
		prefix = "{env.api_sub_path}"
	}
	for path, pathItem := range swagger.Paths.Paths {
		for method, operation := range GetPathOperationMap(pathItem) {
			// 网关的资源路径为 basePath + path
			routePath := joinPaths(swagger.BasePath, path)
			route := routeMap[routePath+":"+method]
			if route != nil && operation.ID == "" {
				operation.ID = route.OperationID
			}
			if route != nil && route.Config != nil {
				// the route configs are shared, copy it to keep merging idempotent
				config := *route.Config
				config.Backend.Method = cmp.Or(config.Backend.Method, method)
				config.Backend.Path = cmp.Or(config.Backend.Path, routePath)
				// 其他后端不是当前应用，不加子路径
				if config.Backend.Name == "" || config.Backend.Name == "default" {
					config.Backend.Path = prefix + config.Backend.Path
				}
				config.Backend.MatchSubpath = config.MatchSubpath
				operation.ID = cmp.Or(config.OperationID, operation.ID)
				operation.AddExtension("x-bk-apigateway-resource", &config)
				continue
			}
			// 没有网关配置的接口只补上后端，保留 swagger 中已有的资源配置
			resource, _ := operation.Extensions["x-bk-apigateway-resource"].(map[string]any)
			if resource["backend"] != nil {
				continue
			}
			if resource == nil {
				resource = map[string]any{}
			}
			resource["backend"] = map[string]any{"method": method, "path": prefix + routePath}
			operation.AddExtension("x-bk-apigateway-resource", resource)
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
		toolMap[tool] = struct{}{}
	}

	// Iterate through all paths in the Swagger specification
	for path, pathItem := range swagger.Paths.Paths {
		operationMap := GetPathOperationMap(pathItem)
		// Process each operation in the operation map
		for method, operation := range operationMap {
			// 如果配置了mcp,则校验
			c := routeMap[joinPaths(swagger.BasePath, path)+":"+method]
			if c != nil && c.Config != nil && c.Config.EnableMcp {
				// 如果指定了operationID,则使用指定的
				if c.Config.OperationID != "" {
					operation.ID = c.Config.OperationID
				}
				// 如果没有指定operationID,则使用自动生成的
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
				log.Fatalf(
					"tool: %s not found in any resource with EnableMcp=true. "+
						"Please ensure the tool is configured with EnableMcp=true in its resource configuration.",
					tool,
				)
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
