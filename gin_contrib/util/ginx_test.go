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

import "testing"

func TestGenerateOperationID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// 标准用例
		{
			name:     "basic case",
			input:    "github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/example/api.UpdateProduct",
			expected: "api_update_product",
		},

		// 包名含大写字母
		{
			name:     "uppercase package",
			input:    "FooPkg.BarHandler",
			expected: "foopkg_bar_handler",
		},

		// 函数名含连续大写
		{
			name:     "consecutive caps",
			input:    "util.ParseJSON",
			expected: "util_parse_json",
		},

		// 函数名含数字
		{
			name:     "alphanumeric",
			input:    "model.UserV2",
			expected: "model_user_v2",
		},

		// 无包名场景
		{
			name:     "no package",
			input:    "StandaloneFunc",
			expected: "_standalone_func", // 注意前导下划线
		},

		// 多级路径
		{
			name:     "nested path",
			input:    "foo/bar/baz.Quux",
			expected: "baz_quux",
		},

		// 保留全大写缩写
		{
			name:     "uppercase acronym",
			input:    "network.HTTPHandler",
			expected: "network_http_handler",
		},

		// 含版本号
		{
			name:     "versioned",
			input:    "v2.GetData",
			expected: "v2_get_data",
		},

		// 复杂驼峰
		{
			name:     "complex camelcase",
			input:    "manager.InitVIPConfig",
			expected: "manager_init_vip_config",
		},
		{
			name:     "complex camelcase with numbers",
			input:    "github.com/gin-gonic/gin.WrapH.func1",
			expected: "gin_wraph_func1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateOperationID(tt.input); got != tt.expected {
				t.Errorf("GenerateOperationID(%q) =\nGot:  %q\nWant: %q", tt.input, got, tt.expected)
			}
		})
	}
}
