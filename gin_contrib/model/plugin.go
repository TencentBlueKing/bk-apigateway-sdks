package model

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	validator "github.com/go-playground/validator/v10"
	yaml "gopkg.in/yaml.v3"
)

type PluginType string

const (
	PluginTypeBKCors            PluginType = "bk-cors"
	PluginTypeHeaderRewrite     PluginType = "bk-header-rewrite"
	PluginTypeIPRestriction     PluginType = "bk-ip-restriction"
	PluginTypeRateLimit         PluginType = "bk-rate-limit"
	PluginTypeMock              PluginType = "bk-mock"
	PluginTypeAPIBreaker        PluginType = "api-breaker"
	PluginTypeRequestValidation PluginType = "request-validation"
	PluginTypeFaultInjection    PluginType = "fault-injection"
)

type PluginConfig struct {
	Type PluginType `json:"type" yaml:"type"`
	YAML string     `json:"yaml" yaml:"yaml,omitempty"`
}

// CorsConfig #################### CORS 插件 ####################
type CorsConfig struct {
	AllowOrigins        string   `yaml:"allow_origins" validate:"required"`
	AllowMethods        string   `yaml:"allow_methods,omitempty"`
	AllowOriginsByRegex []string `yaml:"allow_origins_by_regex,omitempty"`
	AllowHeaders        string   `yaml:"allow_headers,omitempty"`
	ExposeHeaders       string   `yaml:"expose_headers"`
	MaxAge              int      `yaml:"max_age,omitempty" validate:"min=0"`
	AllowCredential     bool     `yaml:"allow_credential,omitempty"`
}

// HeaderRewriteConfig #################### Header 转换插件 ####################
type HeaderRewriteConfig struct {
	Set    []HeaderRewriteValue `yaml:"set,omitempty" validate:"omitempty"`
	Remove []HeaderRewriteValue `yaml:"remove,omitempty" validate:"omitempty"`
}

type HeaderRewriteValue struct {
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}

// IPRestrictionConfig #################### IP 限制插件 ####################
type IPRestrictionConfig struct {
	Whitelist []string `yaml:"whitelist,omitempty" validate:"omitempty,dive,ipv4|cidr|ipv6"`
	Blacklist []string `yaml:"blacklist,omitempty" validate:"omitempty,dive,ipv4|cidr|ipv6"`
	Message   string   `yaml:"message,omitempty" validate:"omitempty,min=1,max=1024"`
}

// RatePolicy #################### 频率控制插件 ####################
type RatePolicy struct {
	Period int `yaml:"period" validate:"required,min=1"`
	Tokens int `yaml:"tokens" validate:"required,min=1"`
}

type RateLimitConfig struct {
	Rates struct {
		Default  []RatePolicy `yaml:"__default" validate:"required"`
		Specials []struct {
			BKAppCode string       `yaml:"bk_app_code" validate:"required,matches=^[a-z][a-z0-9_-]{0,31}$"`
			Policies  []RatePolicy `yaml:"policies" validate:"required"`
		} `yaml:"specials,omitempty"`
	} `yaml:"rates" validate:"required"`
}

// MockConfig #################### Mock 插件 ####################
type MockConfig struct {
	ResponseStatus  int               `yaml:"response_status" validate:"min=100"`
	ResponseExample string            `yaml:"response_example,omitempty"`
	ResponseHeaders map[string]string `yaml:"response_headers,omitempty"`
}

// APIBreakerConfig #################### 熔断插件 ####################
type APIBreakerConfig struct {
	BreakResponseCode    int    `yaml:"break_response_code" validate:"required,min=200,max=599"`
	BreakResponseBody    string `yaml:"break_response_body,omitempty"`
	BreakResponseHeaders []struct {
		Key   string `yaml:"key" validate:"required"`
		Value string `yaml:"value" validate:"required"`
	} `yaml:"break_response_headers,omitempty"`
	MaxBreakerSec int `yaml:"max_breaker_sec" validate:"min=3"`
	Unhealthy     struct {
		HTTPStatuses []int `yaml:"http_statuses" validate:"required,min=1,dive,min=500,max=599"`
		Failures     int   `yaml:"failures" validate:"required,min=1"`
	} `yaml:"unhealthy" validate:"required"`
	Healthy struct {
		HTTPStatuses []int `yaml:"http_statuses" validate:"required,min=1,dive,min=200,max=499"`
		Successes    int   `yaml:"successes" validate:"required,min=1"`
	} `yaml:"healthy" validate:"required"`
}

// RequestValidationConfig #################### 请求校验插件 ####################
type RequestValidationConfig struct {
	BodySchema   map[string]interface{} `yaml:"body_schema,omitempty"`
	HeaderSchema map[string]interface{} `yaml:"header_schema,omitempty"`
	RejectedCode int                    `yaml:"rejected_code" validate:"min=200,max=599"`
	RejectedMsg  string                 `yaml:"rejected_msg,omitempty" validate:"omitempty,max=256"`
}

// FaultInjectionConfig #################### 故障注入插件 ####################
type FaultInjectionConfig struct {
	Abort *struct {
		HTTPStatus int    `yaml:"http_status" validate:"required,min=200"`
		Body       string `yaml:"body,omitempty"`
		Percentage int    `yaml:"percentage,omitempty" validate:"omitempty,min=0,max=100"`
		Vars       string `yaml:"vars,omitempty"`
	} `yaml:"abort,omitempty"`
	Delay *struct {
		Duration   int    `yaml:"duration" validate:"required,min=0"`
		Percentage int    `yaml:"percentage,omitempty" validate:"omitempty,min=0,max=100"`
		Vars       string `yaml:"vars,omitempty"`
	} `yaml:"delay,omitempty"`
}

// BuildResourcePluginConfigWithType 创建插件配置（带自动校验和格式处理）
func BuildResourcePluginConfigWithType(pluginType PluginType, config interface{}) *PluginConfig {
	// 配置校验
	if err := validateConfig(config); err != nil {
		log.Printf("plugin config validate error: %v\n", err)
		return nil
	}

	// YAML 序列化
	data, err := yaml.Marshal(config)
	if err != nil {
		log.Printf("plugin config marshal error: %v", err)
		return nil
	}

	return &PluginConfig{
		Type: pluginType,
		YAML: strings.ReplaceAll(string(data), `""`, `''`), // 处理空字符串
	}
}

func BuildStagePluginConfigWithType(pluginType PluginType, config interface{}) *PluginConfig {
	// 配置校验
	if err := validateConfig(config); err != nil {
		log.Printf("plugin config validate error: %v\n", err)
		return nil
	}

	// YAML 序列化
	data, err := yaml.Marshal(config)
	if err != nil {
		log.Printf("plugin config marshal error: %v", err)
		return nil
	}

	return &PluginConfig{
		Type: pluginType,
		YAML: string(data),
	}
}

// validateConfig 配置校验
func validateConfig(config interface{}) error {
	if reflect.ValueOf(config).IsZero() {
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("config validate error: %w", err)
	}
	return nil
}
