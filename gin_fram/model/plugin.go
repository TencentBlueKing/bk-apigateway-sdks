package model

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type PluginType string

const (
	PluginTypeBKCors PluginType = "bk-cors"
)

type PluginConfig struct {
	Type PluginType `json:"type" yaml:"type"`
	YAML string     `json:"yaml" yaml:"yaml,omitempty"`
}

// CorsConfig ----------------------
// CORS 插件配置
// ----------------------
type CorsConfig struct {
	AllowOrigins        string   `yaml:"allow_origins" validate:"required"`
	AllowMethods        string   `yaml:"allow_methods,omitempty"`
	AllowOriginsByRegex []string `yaml:"allow_origins_by_regex,omitempty"`
	AllowHeaders        string   `yaml:"allow_headers,omitempty"`
	ExposeHeaders       string   `yaml:"expose_headers"`
	MaxAge              int      `yaml:"max_age,omitempty" validate:"min=0"`
	AllowCredential     bool     `yaml:"allow_credential,omitempty"`
}

// RatePolicy ----------------------
// 频率限制插件配置
// ----------------------
type RatePolicy struct {
	Period int `yaml:"period" validate:"required,min=1"`
	Tokens int `yaml:"tokens" validate:"required,min=1"`
}
type RateLimitConfig struct {
	Rates struct {
		Default []RatePolicy `yaml:"__default"`
	} `yaml:"rates"`
	Specials []struct {
		BKAppCode string       `yaml:"bk_app_code"`
		Policies  []RatePolicy `yaml:"policies"`
	} `yaml:"specials,omitempty"`
}

// BuildPluginConfigWithType
// 核心工具函数
// BuildPluginConfig 创建插件配置（带自动校验和格式处理）
func BuildPluginConfigWithType(pluginType PluginType, config interface{}) *PluginConfig {
	// 配置校验
	if err := validateConfig(config); err != nil {
		log.Fatalf("配置校验失败: %v", err)
		return nil
	}
	// YAML 序列化
	data, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("YAML序列化失败: %v", err)
		return nil
	}
	return &PluginConfig{
		Type: pluginType,
		YAML: strings.ReplaceAll(string(data), `""`, `''`),
	}
}

// validateConfig 配置校验
func validateConfig(config interface{}) error {
	// 跳过空配置校验
	if reflect.ValueOf(config).IsZero() {
		return nil
	}
	if err := validator.New().Struct(config); err != nil {
		return fmt.Errorf("配置校验错误: %w", err)
	}
	return nil
}
