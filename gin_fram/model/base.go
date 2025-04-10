package model

// 完整配置结构体定义
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
}
type GatewayConfig struct {
	Description   string
	DescriptionEn string
	IsPublic      bool
	APIType       string
	Maintainers   []string
}
type StageConfig struct {
	Name           string
	Description    string
	DescriptionEn  string
	BackendSubPath string
	BackendTimeout int
	BackendHost    string
	PluginConfigs  []*PluginConfig
}
type GrantPermissionConfig struct {
	GatewayApps  []string
	ResourceApps map[string][]string
}
type ResourceDocConfig struct {
	BaseDir string `validate:"required,startswith=/"` // 必须为绝对路径
}
