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
	"testing"
)

func TestConvertExpressPathToSwagger(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "基础参数转换",
			input:    "/api/:id",
			expected: "/api/{id}",
		},
		{
			name:     "多个参数混合大小写",
			input:    "/users/:userId/orders/:orderId",
			expected: "/users/{userId}/orders/{orderId}",
		},
		{
			name:     "包含特殊字符的参数名",
			input:    "/data/:data-id/:data_name",
			expected: "/data/{data-id}/{data_name}",
		},
		{
			name:     "根路径参数",
			input:    "/:version",
			expected: "/{version}",
		},
		{
			name:     "路径末尾参数",
			input:    "/posts/:post-id/",
			expected: "/posts/{post-id}/",
		},
		{
			name:     "无参数路径",
			input:    "/static/path/to/resource",
			expected: "/static/path/to/resource",
		},
		{
			name:     "空字符串处理",
			input:    "",
			expected: "",
		},
		{
			name:     "无效参数格式处理",
			input:    "/api/:",
			expected: "/api/:",
		},
		{
			name:     "多斜杠路径处理",
			input:    "//api///:id",
			expected: "//api///{id}",
		},
		{
			name:     "中文参数支持",
			input:    "/用户/:用户ID/订单/:订单号",
			expected: "/用户/{用户ID}/订单/{订单号}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertExpressPathToSwagger(tt.input)
			if result != tt.expected {
				t.Errorf("输入: %q\n期望: %q\n得到: %q", tt.input, tt.expected, result)
			}
		})
	}
}
