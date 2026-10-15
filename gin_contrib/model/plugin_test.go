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

package model

import "testing"

func TestPlugin(t *testing.T) {
	// CORS 插件配置示例
	corsConfig := CorsConfig{
		AllowOrigins:    "https://example.com",
		AllowMethods:    "GET,POST",
		AllowHeaders:    "X-Custom-Header",
		MaxAge:          3600,
		AllowCredential: true,
	}
	if plugin := BuildResourcePluginConfigWithType(PluginTypeBKCors, corsConfig); plugin.Type != PluginTypeBKCors {
		t.Errorf("unexpected plugin: %+v", plugin)
	}
	// 熔断插件配置示例
	breakerConfig := APIBreakerConfig{
		BreakResponseCode: 503,
		MaxBreakerSec:     300,
		Unhealthy: struct {
			HTTPStatuses []int `yaml:"http_statuses" validate:"required,min=1,dive,min=500,max=599"`
			Failures     int   `yaml:"failures" validate:"required,min=1"`
		}{
			HTTPStatuses: []int{500, 502},
			Failures:     3,
		},
		Healthy: struct {
			HTTPStatuses []int `yaml:"http_statuses" validate:"required,min=1,dive,min=200,max=499"`
			Successes    int   `yaml:"successes" validate:"required,min=1"`
		}{
			HTTPStatuses: []int{200},
			Successes:    3,
		},
	}
	if plugin := BuildResourcePluginConfigWithType(PluginTypeAPIBreaker, breakerConfig); plugin.YAML == "" {
		t.Errorf("unexpected plugin: %+v", plugin)
	}
}

func TestInvalidPluginConfig(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("an invalid plugin config should panic")
		}
	}()
	// allow_origins is required
	BuildResourcePluginConfigWithType(PluginTypeBKCors, CorsConfig{MaxAge: 1})
}
