package gen

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/model"
)

//go:embed definition.tpl
var configTemplate []byte

func GenDefinitionYaml(config *model.APIConfig) string {
	// 创建模板
	tmpl, err := template.New("config").Funcs(template.FuncMap{
		"indent": func(n int, s string) string {
			pad := strings.Repeat(" ", n)
			return pad + strings.ReplaceAll(s, "\n", "\n"+pad)
		},
	}).Parse(string(configTemplate))
	if err != nil {
		panic(fmt.Sprintf("gen definition yaml error: %v", err))
	}
	// 渲染输出
	var result strings.Builder
	if err := tmpl.Execute(&result, config); err != nil {
		panic(fmt.Sprintf("gen definition yaml error: %v", err))
	}
	fmt.Println("gen definition yaml:\n" + result.String())
	return result.String()
}
