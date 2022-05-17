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

package apigateway

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
)

//  添加网关关联应用
func (c *Client) AddRelatedApps(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "add_related_apps",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/related-apps/",
	}, opts...)
}

//  申请网关API访问权限
func (c *Client) ApplyPermissions(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "apply_permissions",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/permissions/apply/",
	}, opts...)
}

//  创建资源版本
func (c *Client) CreateResourceVersion(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "create_resource_version",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/resource_versions/",
	}, opts...)
}

//  生成 SDK
func (c *Client) GenerateSdk(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "generate_sdk",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/sdk/",
	}, opts...)
}

//  获取网关公钥
func (c *Client) GetApigwPublicKey(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_apigw_public_key",
		Method: "GET",
		Path:   "/api/v1/apis/{api_name}/public_key/",
	}, opts...)
}

//  获取网关最新版本
func (c *Client) GetLatestResourceVersion(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_latest_resource_version",
		Method: "GET",
		Path:   "/api/v1/apis/{api_name}/resource_versions/latest/",
	}, opts...)
}

//  获取微网关应用权限
func (c *Client) GetMicroGatewayAppPermissions(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_micro_gateway_app_permissions",
		Method: "GET",
		Path:   "/api/v1/edge-controller/micro-gateway/{instance_id}/permissions/",
	}, opts...)
}

//  网关为应用主动授权
func (c *Client) GrantPermissions(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "grant_permissions",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/permissions/grant/",
	}, opts...)
}

//  通过文档归档文件导入资源文档
func (c *Client) ImportResourceDocsByArchive(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "import_resource_docs_by_archive",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/resource-docs/import/by-archive/",
	}, opts...)
}

// 通过 Swagger 格式导入文档 :
func (c *Client) ImportResourceDocsBySwagger(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "import_resource_docs_by_swagger",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/resource-docs/import/by-swagger/",
	}, opts...)
}

//  发布版本
func (c *Client) Release(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "release",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/resource_versions/release/",
	}, opts...)
}

//  回收应用访问网关 API 的权限
func (c *Client) RevokePermissions(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "revoke_permissions",
		Method: "DELETE",
		Path:   "/api/v1/apis/{api_name}/permissions/revoke/",
	}, opts...)
}

//  同步策略
func (c *Client) SyncAccessStrategy(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "sync_access_strategy",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/access_strategies/sync/",
	}, opts...)
}

//  同步网关
func (c *Client) SyncApi(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "sync_api",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/sync/",
	}, opts...)
}

//  同步资源
func (c *Client) SyncResources(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "sync_resources",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/resources/sync/",
	}, opts...)
}

//  同步环境
func (c *Client) SyncStage(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "sync_stage",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/stages/sync/",
	}, opts...)
}

//  更新微网关实例状态
func (c *Client) UpdateMicroGatewayStatus(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "update_micro_gateway_status",
		Method: "PUT",
		Path:   "/api/v1/edge-controller/micro-gateway/{instance_id}/status/",
	}, opts...)
}
