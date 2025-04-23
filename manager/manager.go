/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2017 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package manager

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/TencentBlueKing/bk-apigateway-sdks/apigateway"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/util"
)

const (
	apiGatewayNamespace   = "apigateway"
	stagesNamespace       = "stages"
	permissionsNamespace  = "grant_permissions"
	resourceDocsNamespace = "resource_docs"
)

type apiGatewayResult struct {
	Code      int                    `json:"code"`
	HasResult bool                   `json:"result"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
}

// Manager is the manager of apigw, it helps to sync apigw configs and get apigw infomations.
type Manager struct {
	apiName    string
	definition *Definition
	client     *apigateway.Client
	config     *bkapi.ClientConfig
}

func (m *Manager) requestWithBody(
	operation define.Operation,
	body map[string]interface{},
) (map[string]interface{}, error) {
	return m.request(operation.SetBody(body))
}

func (m *Manager) requestWithFile(
	operation define.Operation,
	name string,
	file *os.File,
) (map[string]interface{}, error) {
	return m.request(operation.SetFile(name, file))
}

func (m *Manager) request(operation define.Operation) (map[string]interface{}, error) {
	var result apiGatewayResult
	_, err := operation.
		SetPathParams(map[string]string{
			"api_name": m.apiName,
		}).
		SetResult(&result).
		Request()
	if err != nil {
		return nil, errors.Wrapf(err, "request to %v failed", operation)
	}

	if result.Code == 0 {
		return result.Data, nil
	}

	return result.Data, errors.Wrapf(
		ErrApigatewayRequest,
		"code: %d, message: %s",
		result.Code,
		result.Message,
	)
}

// LoadDefinition will load the definition from the file.
func (m *Manager) LoadDefinition(path string) error {
	rendered, err := ioutil.ReadFile(path)
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

// GetDefinition return the definition.
func (m *Manager) GetDefinition() *Definition {
	return m.definition
}

// GetPublicKey fetch the public key info from apigw.
func (m *Manager) GetPublicKey() (map[string]interface{}, error) {
	return m.request(m.client.GetApigwPublicKey())
}

// GetPublicKey fetch the public key from apigw.
func (m *Manager) GetPublicKeyString() (string, error) {
	info, err := m.GetPublicKey()
	if err != nil {
		return "", err
	}

	value, ok := info["public_key"]
	if !ok {
		return "", errors.Wrapf(ErrApiGatewayPublicKeyNotFound, m.apiName)
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

// GetLatestResourceVersion get the latest resource version from apigw.
func (m *Manager) GetLatestResourceVersion() (map[string]interface{}, error) {
	return m.request(m.client.GetLatestResourceVersion())
}

// SyncBasicInfo sync the basic info from definition under the namespace to apigw.
func (m *Manager) SyncBasicInfo() (map[string]interface{}, error) {
	data, err := m.definition.Get(apiGatewayNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", apiGatewayNamespace)
	}

	return m.requestWithBody(m.client.SyncAPI(), data)
}

// SyncStagesConfig sync the stages config from definition under the namespace to apigw.
func (m *Manager) SyncStagesConfig() (map[string]interface{}, error) {
	stages, err := m.definition.GetArray(stagesNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
	}
	resultMap := make(map[string]interface{})
	for _, stage := range stages {
		result, err := m.requestWithBody(m.client.SyncStage(), stage)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
		}
		resultMap[fmt.Sprintf("%v", stage["name"])] = result
	}
	return resultMap, nil
}

// SyncPluginConfig sync the plugin config from definition under the namespace to apigw.
func (m *Manager) SyncPluginConfig(namespace string) (map[string]interface{}, error) {
	data, err := m.definition.Get(namespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", namespace)
	}

	return m.requestWithBody(m.client.SyncAccessStrategy(), data)
}

// SyncResourcesConfig sync the resources config from definition under the namespace to apigw.
func (m *Manager) SyncResourcesConfig(resources map[string]interface{}) (map[string]interface{}, error) {
	return m.requestWithBody(m.client.SyncResources(), resources)
}

// SyncResourceDocByArchive sync the resource doc from archive to apigw.
func (m *Manager) SyncResourceDocByArchive() (map[string]interface{}, error) {
	data, err := m.definition.Get(resourceDocsNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", resourceDocsNamespace)
	}
	baseDir := data["basedir"].(string)
	err = util.ZipDirectory(baseDir, baseDir+"resources_docs.zip", ".md")
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to zip %s", data["base_dir"].(string))
	}
	// 上传资源文档
	resourceDocsFile, err := os.Open(baseDir + "/resources_docs.zip")
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to read %s", baseDir+"/resources_docs.zip")
	}
	return m.requestWithFile(m.client.ImportResourceDocsByArchive(), "file", resourceDocsFile)
}

// ApplyPermissions apply the permissions under the namespace to apigw.
func (m *Manager) ApplyPermissions(namespace string) (map[string]interface{}, error) {
	data, err := m.definition.Get(namespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", namespace)
	}

	return m.requestWithBody(m.client.ApplyPermissions(), data)
}

// GrantPermissions grant the permissions under the namespace to apigw.
func (m *Manager) GrantPermissions() (map[string]interface{}, error) {
	datas, err := m.definition.GetArray(permissionsNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", permissionsNamespace)
	}
	resultMap := make(map[string]interface{})
	for i, data := range datas {
		param := map[string]interface{}{
			"target_app_code": data["bk_app_code"],
			"grant_dimension": data["grant_dimension"],
		}
		resourceNames, ok := data["resource_names"]
		if ok {
			param["resource_names"] = resourceNames
		}
		result, err := m.requestWithBody(m.client.GrantPermissions(), param)
		if err != nil {
			return nil, err
		}
		resultMap[fmt.Sprintf("result_%d", i)] = result
	}
	return resultMap, nil
}

// CreateResourceVersion create a resource version defined in the namespace.
func (m *Manager) CreateResourceVersion(version string, comment string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"version": version,
		"comment": comment,
	}
	return m.requestWithBody(m.client.CreateResourceVersion(), data)
}

// Release release the resource version defined in the namespace.
func (m *Manager) Release(version string) (map[string]interface{}, error) {
	stages, err := m.definition.GetArray(stagesNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", permissionsNamespace)
	}
	var stageNames []string
	for _, stage := range stages {
		stageNames = append(stageNames, stage["name"].(string))
	}
	data := map[string]interface{}{
		"stage_names": stageNames,
		"version":     version,
	}
	return m.requestWithBody(m.client.Release(), data)
}

// NewManager create a new manager.
func NewManager(
	apiName string,
	config bkapi.ClientConfig,
	definition *Definition,
	clientFactory func(
		configProvider define.ClientConfigProvider, opts ...define.BkApiClientOption,
	) (*apigateway.Client, error),
) (*Manager, error) {
	client, err := clientFactory(config, bkapi.OptJsonBodyProvider(), bkapi.JsonResultProvider())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create apigateway client")
	}

	return &Manager{
		apiName:    apiName,
		config:     &config,
		client:     client,
		definition: definition,
	}, nil
}

// NewDefaultManager create a new default manager.
func NewDefaultManager(apiName string, config bkapi.ClientConfig) (*Manager, error) {
	return NewManager(apiName, config, nil, apigateway.New)
}

// NewManagerFrom file will create a new manager from the file.
func NewManagerFrom(
	apiName string,
	config bkapi.ClientConfig,
	path string,
) (*Manager, error) {
	manager, err := NewDefaultManager(apiName, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create manager")
	}

	return manager, manager.LoadDefinition(path)
}
