/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	gentleman "gopkg.in/h2non/gentleman.v2"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/util"
)

const (
	apiGatewayNamespace       = "apigateway"
	stagesNamespace           = "stages"
	grantPermissionsNamespace = "grant_permissions"
	applyPermissionsNamespace = "apply_permissions"
	resourceDocsNamespace     = "resource_docs"
	relatedAppsNamespace      = "related_apps"
)

// ConfigFromEnv returns the config to call the prod stage of bk-apigateway, see bkapi.ConfigFromEnv.
func ConfigFromEnv() bkapi.Config {
	return bkapi.ConfigFromEnv("bk-apigateway", "prod")
}

// Manager is the manager of apigw, it helps to sync apigw configs and get apigw information.
// It calls the v2 apis of bk-apigateway.
type Manager struct {
	gatewayName string
	appCode     string
	definition  *Definition
	client      *gentleman.Client
}

// send sends the request to bk-apigateway, and returns the data of the response, which is like {"data": ...}.
//
// A non-2xx response, whose body is like {"error": ...}, never reaches the decoding below: the client created
// by bkapi.New turns it into a *bkapi.Error returned by Send.
func send[T any](ctx context.Context, req *gentleman.Request) (T, error) {
	var body struct {
		Data T `json:"data"`
	}
	res, err := req.Use(bkapi.WithContext(ctx)).Send()
	if err != nil {
		return body.Data, err
	}
	// res.JSON takes an empty body as a success
	if err := json.Unmarshal(res.Bytes(), &body); err != nil {
		return body.Data, errors.Wrap(err, "failed to decode response")
	}
	return body.Data, nil
}

// LoadDefinition loads the definition from the file.
func (m *Manager) LoadDefinition(path string) error {
	rendered, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", path)
	}
	definition, err := NewDefinitionFromYaml(rendered)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %s", path)
	}

	m.definition = definition
	return nil
}

// GetDefinition returns the definition.
func (m *Manager) GetDefinition() *Definition {
	return m.definition
}

// GetPublicKey fetches the public key info from apigw.
func (m *Manager) GetPublicKey(ctx context.Context) (map[string]any, error) {
	return send[map[string]any](ctx, m.client.Get().
		AddPath("/api/v2/sync/gateways/:gateway_name/public_key/").
		Param("gateway_name", m.gatewayName))
}

// GetPublicKeyString fetches the public key from apigw.
func (m *Manager) GetPublicKeyString(ctx context.Context) (string, error) {
	info, err := m.GetPublicKey(ctx)
	if err != nil {
		return "", err
	}

	value, ok := info["public_key"]
	if !ok {
		return "", errors.Wrap(ErrApiGatewayPublicKeyNotFound, m.gatewayName)
	}

	publicKey, ok := value.(string)
	if !ok {
		return "", errors.Wrapf(
			ErrApiGatewayPublicKeyTypeNotSupported,
			"expected %T, got %T", publicKey, value,
		)
	}

	return publicKey, nil
}

// GetLatestResourceVersion gets the latest resource version from apigw.
func (m *Manager) GetLatestResourceVersion(ctx context.Context) (map[string]any, error) {
	return send[map[string]any](ctx, m.client.Get().
		AddPath("/api/v2/sync/gateways/:gateway_name/resource_versions/latest/").
		Param("gateway_name", m.gatewayName))
}

// ResourceVersionExists checks whether the resource version exists in apigw.
func (m *Manager) ResourceVersionExists(ctx context.Context, version string) (bool, error) {
	result, err := send[struct {
		Count int `json:"count"`
	}](ctx, m.client.Get().
		AddPath("/api/v2/sync/gateways/:gateway_name/resource_versions/").
		Param("gateway_name", m.gatewayName).
		SetQuery("version", version))
	return result.Count > 0, err
}

// SyncBasicInfo syncs the basic info from definition under the namespace to apigw.
func (m *Manager) SyncBasicInfo(ctx context.Context) (map[string]any, error) {
	data, err := m.definition.Get(apiGatewayNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", apiGatewayNamespace)
	}
	return send[map[string]any](ctx, m.client.Post().
		AddPath("/api/v2/sync/gateways/:gateway_name/").
		Param("gateway_name", m.gatewayName).
		JSON(data))
}

// SyncStagesConfig syncs the stages config from definition under the namespace to apigw.
func (m *Manager) SyncStagesConfig(ctx context.Context) (map[string]any, error) {
	stages, err := m.definition.GetArray(stagesNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
	}
	results := make(map[string]any, len(stages))
	for _, stage := range stages {
		result, err := send[map[string]any](ctx, m.client.Post().
			AddPath("/api/v2/sync/gateways/:gateway_name/stages/").
			Param("gateway_name", m.gatewayName).
			JSON(stage))
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to sync stage %v", stage["name"])
		}
		results[fmt.Sprint(stage["name"])] = result
	}
	return results, nil
}

// SyncStageMcpConfig syncs the mcp servers of the stages from definition to apigw.
// The stages without mcp_servers are skipped, as the api requires it.
func (m *Manager) SyncStageMcpConfig(ctx context.Context) (map[string]any, error) {
	stages, err := m.definition.GetArray(stagesNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
	}
	results := make(map[string]any, len(stages))
	for _, stage := range stages {
		if _, ok := stage["mcp_servers"]; !ok {
			continue
		}
		name := fmt.Sprint(stage["name"])
		result, err := send[[]map[string]any](ctx, m.client.Post().
			AddPath("/api/v2/sync/gateways/:gateway_name/stages/:stage_name/mcp-servers/").
			Param("gateway_name", m.gatewayName).
			Param("stage_name", name).
			JSON(stage))
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to sync the mcp servers of stage %s", name)
		}
		results[name] = result
	}
	return results, nil
}

// SyncResourcesConfig syncs the resources config to apigw.
func (m *Manager) SyncResourcesConfig(ctx context.Context, resources map[string]any) (map[string]any, error) {
	return send[map[string]any](ctx, m.client.Post().
		AddPath("/api/v2/sync/gateways/:gateway_name/resources/").
		Param("gateway_name", m.gatewayName).
		JSON(resources))
}

// SyncResourceDocByArchive zips the markdown docs under resource_docs.basedir, and uploads them to apigw.
func (m *Manager) SyncResourceDocByArchive(ctx context.Context) error {
	data, err := m.definition.Get(resourceDocsNamespace)
	if err != nil {
		return errors.WithMessagef(err, "failed to get %s", resourceDocsNamespace)
	}
	baseDir, _ := data["basedir"].(string)
	if baseDir == "" {
		return errors.Errorf("%s.basedir is required", resourceDocsNamespace)
	}
	zipPath := filepath.Join(baseDir, "resources_docs.zip")
	err = util.ZipDirectory(baseDir, zipPath, ".md")
	if err != nil {
		return errors.WithMessagef(err, "failed to zip %s", baseDir)
	}
	// 上传资源文档
	archive, err := os.Open(zipPath)
	if err != nil {
		return errors.WithMessagef(err, "failed to read %s", zipPath)
	}
	defer archive.Close()

	_, err = send[any](ctx, m.client.Post().
		AddPath("/api/v2/sync/gateways/:gateway_name/resource-docs/").
		Param("gateway_name", m.gatewayName).
		File("file", archive))
	return err
}

// targetGateway removes the gateway_name (or the legacy api_name, which takes precedence) from the permission,
// which is the gateway to apply for or grant, and defaults to the managed gateway.
func (m *Manager) targetGateway(permission map[string]any) string {
	name, ok := permission["api_name"]
	if !ok {
		name = permission["gateway_name"]
	}
	delete(permission, "api_name")
	delete(permission, "gateway_name")
	if name == nil {
		return m.gatewayName
	}
	return fmt.Sprint(name)
}

// ApplyPermissions applies for the permissions defined in apply_permissions, which can be of other gateways.
func (m *Manager) ApplyPermissions(ctx context.Context) (map[string]any, error) {
	if m.definition.definition[applyPermissionsNamespace] == nil {
		return nil, nil
	}
	permissions, err := m.definition.GetArray(applyPermissionsNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", applyPermissionsNamespace)
	}

	results := make(map[string]any, len(permissions))
	for i, definedPermission := range permissions {
		permission := maps.Clone(definedPermission)
		if _, ok := permission["target_app_code"]; !ok {
			permission["target_app_code"] = m.appCode
		}
		if _, ok := permission["applicant"]; !ok {
			permission["applicant"] = permission["target_app_code"]
		}
		if permission["grant_dimension"] == nil {
			permission["grant_dimension"] = "gateway"
		}

		result, err := send[map[string]any](ctx, m.client.Post().
			AddPath("/api/v2/open/gateways/:gateway_name/permissions/apply/").
			Param("gateway_name", m.targetGateway(permission)).
			JSON(permission))
		if err != nil {
			return nil, err
		}
		results[fmt.Sprintf("result_%d", i)] = result
	}
	return results, nil
}

// GrantPermissions grants the permissions defined in grant_permissions to apps.
func (m *Manager) GrantPermissions(ctx context.Context) error {
	if m.definition.definition[grantPermissionsNamespace] == nil {
		return nil
	}
	permissions, err := m.definition.GetArray(grantPermissionsNamespace)
	if err != nil {
		return errors.WithMessagef(err, "failed to get %s", grantPermissionsNamespace)
	}

	for _, definedPermission := range permissions {
		permission := maps.Clone(definedPermission)
		// v1 的配置用 bk_app_code 表示被授权的应用，用 api 表示按网关授权
		if _, ok := permission["target_app_code"]; !ok {
			permission["target_app_code"] = permission["bk_app_code"]
		}
		delete(permission, "bk_app_code")
		if dimension := permission["grant_dimension"]; dimension == nil || dimension == "api" {
			permission["grant_dimension"] = "gateway"
		}
		_, err := send[any](ctx, m.client.Post().
			AddPath("/api/v2/sync/gateways/:gateway_name/permissions/grant/").
			Param("gateway_name", m.targetGateway(permission)).
			JSON(permission))
		if err != nil {
			return err
		}
	}
	return nil
}

// AddRelatedApps adds the apps defined in related_apps as the related apps of the gateway.
func (m *Manager) AddRelatedApps(ctx context.Context) error {
	// related_apps 可能是 YAML 解析出的 []any，也可能是代码里构造的 []string
	var relatedApps []string
	raw, err := json.Marshal(m.definition.definition[relatedAppsNamespace])
	if err == nil {
		err = json.Unmarshal(raw, &relatedApps)
	}
	if err != nil {
		return errors.Wrapf(err, "%s must be a list of app codes", relatedAppsNamespace)
	}
	if len(relatedApps) == 0 {
		return nil
	}

	_, err = send[any](ctx, m.client.Post().
		AddPath("/api/v2/sync/gateways/:gateway_name/related-apps/").
		Param("gateway_name", m.gatewayName).
		JSON(map[string]any{"related_app_codes": relatedApps}))
	return err
}

// CreateResourceVersion creates a resource version.
func (m *Manager) CreateResourceVersion(ctx context.Context, version, comment string) (map[string]any, error) {
	data := map[string]any{
		"version": version,
		"comment": comment,
	}
	return send[map[string]any](ctx, m.client.Post().
		AddPath("/api/v2/sync/gateways/:gateway_name/resource_versions/").
		Param("gateway_name", m.gatewayName).
		JSON(data))
}

// Release releases the resource version to the stages defined in the definition.
func (m *Manager) Release(ctx context.Context, version, comment string) (map[string]any, error) {
	stages, err := m.definition.GetArray(stagesNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
	}
	stageNames := make([]string, 0, len(stages))
	for _, stage := range stages {
		stageNames = append(stageNames, fmt.Sprint(stage["name"]))
	}
	data := map[string]any{
		"stage_names": stageNames,
		"version":     version,
		"comment":     comment,
	}
	return send[map[string]any](ctx, m.client.Post().
		AddPath("/api/v2/sync/gateways/:gateway_name/resource_versions/release/").
		Param("gateway_name", m.gatewayName).
		JSON(data))
}

// NewManager creates a new manager, the config is used to call bk-apigateway, see ConfigFromEnv.
func NewManager(gatewayName string, config bkapi.Config, definition *Definition) (*Manager, error) {
	client, err := bkapi.New(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bkapi client")
	}

	return &Manager{
		gatewayName: gatewayName,
		appCode:     config.AppCode,
		definition:  definition,
		client:      client,
	}, nil
}

// NewDefaultManager creates a new manager without a definition.
func NewDefaultManager(gatewayName string, config bkapi.Config) (*Manager, error) {
	return NewManager(gatewayName, config, nil)
}

// NewManagerFrom creates a new manager with the definition loaded from the file.
func NewManagerFrom(gatewayName string, config bkapi.Config, path string) (*Manager, error) {
	manager, err := NewDefaultManager(gatewayName, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create manager")
	}

	return manager, manager.LoadDefinition(path)
}
