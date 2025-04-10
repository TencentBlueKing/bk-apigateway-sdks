# reference: https://github.com/TencentBlueKing/bkpaas-python-sdk/tree/master/sdks/apigw-manager
spec_version: 2
release:
  version: "{{.Release.Version}}"
  title: "{{.Release.Title}}"
  comment: "{{.Release.Comment}}"
apigateway:
  description: "{{.APIGateway.Description}}"
  description_en: "{{.APIGateway.DescriptionEn}}"
  is_public: {{.APIGateway.IsPublic}}
  api_type: {{.APIGateway.APIType}}
  maintainers:{{if .APIGateway.Maintainers}}
  {{- range .APIGateway.Maintainers}}
  - "{{.}}"
  {{- end}}
  {{else}}
  - "admin"
  {{end}}
stages:
  - name: "{{.Stage.Name}}"
    description: "{{.Stage.Description}}"
    description_en: "{{.Stage.DescriptionEn}}"
    {{- if .Stage.BackendSubPath}}
    vars:
      api_sub_path: {{.Stage.BackendSubPath}}
    {{- else}}
    vars: {}
    {{- end}}
    backends:
      - name: "default"
        config:
          timeout: {{.Stage.BackendTimeout}}
          loadbalance: "roundrobin"
          hosts:
            - host: "{{.Stage.BackendHost}}"
              weight: 100
    {{- if .Stage.PluginConfigs}}
    plugin_configs:
      {{- range .Stage.PluginConfigs}}
      - type: {{.Type}}
        yaml: |-
  {{.ConfigYAML | indent 8}}  // 关键修改点
  {{- end}}
  {{- end}}
  {{- if or .GrantPermissions.GatewayApps .GrantPermissions.ResourceApps}}
grant_permissions:
  {{- range .GrantPermissions.GatewayApps}}
  - bk_app_code: "{{.}}"
    grant_dimension: "gateway"
  {{- end}}
  {{- range $app_code, $resource_names := .GrantPermissions.ResourceApps}}
  - bk_app_code: "{{$app_code}}"
    grant_dimension: "resource"
    resource_names:
      {{- range $resource_names}}
      - "{{.}}"
      {{- end}}
  {{- end}}
  {{- end}}
related_apps:
  {{- if .RelatedApps}}
- "{{index .RelatedApps 0}}"
  {{- else}}
- ""
  {{- end}}
resource_docs:
  {{- if .ResourceDocs.BaseDir}}
basedir: "{{.ResourceDocs.BaseDir}}"
  {{- else}}
basedir: ""
  {{- end}}