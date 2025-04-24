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
	BuildResourcePluginConfigWithType(PluginTypeBKCors, corsConfig)
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
	BuildResourcePluginConfigWithType(PluginTypeAPIBreaker, breakerConfig)
}
