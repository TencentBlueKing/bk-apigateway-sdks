package model

type APIConfig struct {
	Release          ReleaseConfig
	APIGateway       GatewayConfig
	Stage            StageConfig
	GrantPermissions GrantPermissionConfig
	RelatedApps      []string
	ResourceDocs     ResourceDocConfig
}

type ReleaseConfig struct {
	Version string
	Title   string
	Comment string
	NoPub   bool // 是否不发布
}

type GatewayConfig struct {
	Description   string
	DescriptionEn string
	IsPublic      bool
	APIType       string
	Maintainers   []string
}

type StageConfig struct {
	Name           string          // 环境名称
	Description    string          // 环境描述
	DescriptionEn  string          // 环境描述英文
	BackendSubPath string          // 后端服务前缀路径
	BackendTimeout int             // 后端服务超时时间
	BackendHost    string          // 后端服务地址
	PluginConfigs  []*PluginConfig // 插件配置
}

type GrantPermissionConfig struct {
	GatewayApps  []string
	ResourceApps map[string][]string
}

type ResourceDocConfig struct {
	BaseDir  string `validate:"required,startswith=/"` // 必须为绝对路径
	Language string
}

// APIGatewayResourceConfig resource 配置结构体定义
type APIGatewayResourceConfig struct {
	ResourceBasicConfig                 // 资源基础配置
	Backend             BackendConfig   `json:"backend" yaml:"backend"`                                 // 后端配置
	PluginConfigs       []*PluginConfig `json:"pluginConfigs,omitempty" yaml:"pluginConfigs,omitempty"` // 插件配置
	AuthConfig          AuthConfig      `json:"authConfig" yaml:"authConfig"`                           // 认证配置
}

type ResourceBasicConfig struct {
	IsPublic             bool `json:"isPublic" yaml:"isPublic"`                         // 是否公开
	AllowApplyPermission bool `json:"allowApplyPermission" yaml:"allowApplyPermission"` // 是否允许申请权限
	MatchSubpath         bool `json:"matchSubpath" yaml:"matchSubpath"`                 // 是否匹配子路径
	EnableWebsocket      bool `json:"enableWebsocket" yaml:"enableWebsocket"`           // 是否启用 websocket
}

type Option func(*APIGatewayResourceConfig)

func NewAPIGatewayResourceConfig(config ResourceBasicConfig, opts ...Option) APIGatewayResourceConfig {
	cfg := &APIGatewayResourceConfig{
		ResourceBasicConfig: config,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return *cfg
}

// WithBackend 设置资源后端配置
func (c *ResourceBasicConfig) WithBackend(backend BackendConfig) Option {
	return func(config *APIGatewayResourceConfig) {
		config.Backend = backend
	}
}

// WithPublic 设置资源公开配置
func (c *ResourceBasicConfig) WithPublic(public bool) Option {
	return func(config *APIGatewayResourceConfig) {
		config.IsPublic = public
	}
}

// WithAllowApplyPermission 设置资源申请权限配置
func (c *ResourceBasicConfig) WithAllowApplyPermission(allow bool) Option {
	return func(config *APIGatewayResourceConfig) {
		config.AllowApplyPermission = allow
	}
}

// WithMatchSubpath 设置资源匹配子路径配置
func (c *ResourceBasicConfig) WithMatchSubpath(matchSubpath bool) Option {
	return func(config *APIGatewayResourceConfig) {
		config.MatchSubpath = matchSubpath
	}
}

// WithEnableWebsocket 设置资源启用 websocket 配置
func (c *ResourceBasicConfig) WithEnableWebsocket(enableWebsocket bool) Option {
	return func(config *APIGatewayResourceConfig) {
		config.EnableWebsocket = enableWebsocket
	}
}

// WithPluginConfig 设置资源插件配置
func (c *ResourceBasicConfig) WithPluginConfig(pluginConfigs ...*PluginConfig) Option {
	return func(config *APIGatewayResourceConfig) {
		// 需要校验插件是否重复
		for _, pluginConfig := range pluginConfigs {
			for _, existPluginConfig := range config.PluginConfigs {
				if existPluginConfig.Type == pluginConfig.Type {
					panic("plugin type: " + string(pluginConfig.Type) + " is duplicated")
				}
			}
		}
		config.PluginConfigs = append(config.PluginConfigs, pluginConfigs...)
	}
}

// WithAuthConfig 设置资源认证配置
func (c *ResourceBasicConfig) WithAuthConfig(authConfig AuthConfig) Option {
	return func(config *APIGatewayResourceConfig) {
		config.AuthConfig = authConfig
	}
}

// AuthConfig 资源认证配置
type AuthConfig struct {
	UserVerifiedRequired       bool `json:"userVerifiedRequired" yaml:"userVerifiedRequired"`             // 是否需要用户认证
	AppVerifiedRequired        bool `json:"appVerifiedRequired" yaml:"appVerifiedRequired"`               // 是否需要应用认证
	ResourcePermissionRequired bool `json:"resourcePermissionRequired" yaml:"resourcePermissionRequired"` // 是否需要资源权限认证
}

// BackendConfig 资源后端配置
type BackendConfig struct {
	Name         string `json:"name,omitempty" yaml:"name,omitempty"`       // 资源后端服务名称 默认为 default
	Method       string `json:"method,omitempty" yaml:"method,omitempty"`   // 资源后端服务请求方法
	Path         string `json:"path" yaml:"path"`                           // 资源后端服务路径
	MatchSubpath bool   `json:"matchSubpath" yaml:"matchSubpath"`           // 资源是否匹配子路径
	Timeout      int    `json:"timeout,omitempty" yaml:"timeout,omitempty"` // 资源服务超时时间
}
