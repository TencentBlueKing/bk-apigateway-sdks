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
