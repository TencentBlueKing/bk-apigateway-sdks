/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2017-2021 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package manager

import (
	"os"
	"strings"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/flosch/pongo2/v5"
)

// DefintionContext for defintion template engine
type DefintionContext struct {
	apiName string
	config  *bkapi.ClientConfig
}

func (c *DefintionContext) settings() map[string]interface{} {
	return map[string]interface{}{
		"BK_APIGW_NAME": c.apiName,
		"BK_APP_CODE":   c.config.AppCode,
		"BK_APP_SECRET": c.config.AppSecret,
	}
}

func (c *DefintionContext) environ() map[string]string {
	envs := os.Environ()
	envMap := make(map[string]string, len(envs))
	for _, env := range envs {
		kv := strings.SplitN(env, "=", 2)
		switch len(kv) {
		case 0:
			continue
		case 1:
			envMap[kv[0]] = ""
		default:
			envMap[kv[0]] = kv[1]
		}
	}

	return envMap
}

// Context return pongo2 context
func (c *DefintionContext) Context(data interface{}) pongo2.Context {
	settings := c.settings()
	environ := c.environ()

	return pongo2.Context{
		"settings": settings,
		"environ":  environ,
		"data":     data,
	}
}

// NewDefinitionContext return new defintion context
func NewDefinitionContext(apiName string, config *bkapi.ClientConfig) *DefintionContext {
	return &DefintionContext{
		apiName: apiName,
		config:  config,
	}
}
