package util

import (
	"encoding/json"
	"strings"

	yaml "gopkg.in/yaml.v2"
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
	segments := strings.Split(path, "/")
	for i, seg := range segments {
		if strings.HasPrefix(seg, ":") && len(seg) > 1 {
			segments[i] = "{" + seg[1:] + "}"
		}
	}
	return strings.Join(segments, "/")
}
