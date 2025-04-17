package util

import (
	"encoding/json"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

func JsonToYAML(jsonData []byte) ([]byte, error) {
	var jsonObj interface{}
	if err := json.Unmarshal(jsonData, &jsonObj); err != nil {
		return nil, err
	}
	return yaml.Marshal(jsonObj)
}

// ConvertExpressPathToSwagger 将 /api/:id 转换为 Swagger 风格 的路径 /api/{id}
func ConvertExpressPathToSwagger(path string) string {
	// 正则匹配 : 开头的参数（如 :id, :version）
	re := regexp.MustCompile(`:(\w+)`)
	// 替换为 {param} 格式
	swaggerPath := re.ReplaceAllStringFunc(path, func(match string) string {
		paramName := strings.TrimPrefix(match, ":")
		return "{" + paramName + "}"
	})
	return swaggerPath
}
