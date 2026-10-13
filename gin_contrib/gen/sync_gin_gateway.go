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

package gen

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/model"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/manager"
)

// SyncGinGateway syncs the gateway defined in definition.yaml and resources.yaml under baseDir to bk-apigateway,
// creates a resource version and releases it. deleteResources deletes the resources not in resources.yaml.
//
// It is for the sync command of an app, and exits the process if any step fails.
func SyncGinGateway(baseDir, gatewayName string, config *model.APIConfig, deleteResources bool) {
	ctx := context.Background()
	defaultManager, err := manager.NewManagerFrom(
		gatewayName,
		manager.ConfigFromEnv(),
		filepath.Join(baseDir, "definition.yaml"),
	)
	if err != nil {
		log.Fatalf("create manager: %v", err)
	}

	// 同步网关基础信息
	info, err := defaultManager.SyncBasicInfo(ctx)
	if err != nil {
		log.Fatalf("sync gateway basic info: %v", err)
	}
	log.Printf("syncing gateway basic info success, info:%v\n", info)

	// 添加网关关联应用，只会新增，不会移除已有的关联应用
	if err := defaultManager.AddRelatedApps(ctx); err != nil {
		log.Fatalf("add related apps: %v", err)
	}

	// 同步网关环境信息
	result, err := defaultManager.SyncStagesConfig(ctx)
	if err != nil {
		log.Fatalf("sync gateway stages: %v", err)
	}
	log.Printf("syncing gateway stage config success, result:%v\n", result)

	// 同步网关资源信息
	resourceFile, err := os.ReadFile(filepath.Join(baseDir, "resources.yaml"))
	if err != nil {
		log.Fatalf("read resources file: %v", err)
	}
	log.Printf("call sync_apigw_resources with resources:%s\n", resourceFile)

	syncResourcesArgs := map[string]any{
		"content": string(resourceFile),
		"delete":  deleteResources,
	}
	if config.ResourceDocs.Language != "" {
		syncResourcesArgs["doc_language"] = config.ResourceDocs.Language
	}
	result, err = defaultManager.SyncResourcesConfig(ctx, syncResourcesArgs)
	if err != nil {
		log.Fatalf("sync gateway resources: %v", err)
	}
	log.Printf("syncing gateway resource config success, result:%v\n", result)

	// 同步授权信息
	if err := defaultManager.GrantPermissions(ctx); err != nil {
		log.Fatalf("grant permissions: %v", err)
	}
	log.Printf("granting gateway permissions success\n")

	// 同步资源文档信息
	if config.ResourceDocs.BaseDir != "" {
		if err := defaultManager.SyncResourceDocByArchive(ctx); err != nil {
			log.Fatalf("sync resource docs: %v", err)
		}
		log.Printf("syncing gateway resource doc success\n")
	}

	// 生成资源版本
	// apigw-manager(python) 会持久化资源签名，在版本号一致且资源无变更时复用最新版本；
	// 这里没有可持久化的签名，所以每次都创建新版本
	newVersion := config.Release.Version
	if newVersion == "" {
		newVersion = "0.0.1"
	}
	exists, err := defaultManager.ResourceVersionExists(ctx, newVersion)
	if err != nil {
		log.Fatalf("check resource version %s: %v", newVersion, err)
	}
	if exists {
		publicVersion := strings.SplitN(newVersion, "+", 2)[0]
		newVersion = fmt.Sprintf("%s+%s", publicVersion, time.Now().Format("20060102150405"))
	}
	result, err = defaultManager.CreateResourceVersion(ctx, newVersion, config.Release.Comment)
	if err != nil {
		log.Fatalf("create resource version %s: %v", newVersion, err)
	}
	log.Printf("create gateway resource version success, result:%v\n", result)
	// 发布资源版本
	if !config.Release.NoPub {
		result, err = defaultManager.Release(ctx, newVersion, config.Release.Comment)
		if err != nil {
			log.Fatalf("release resource version %s: %v", newVersion, err)
		}
		log.Printf("release gateway resource version success, result:%v\n", result)
	}

	if config.Stage != nil && config.Stage.EnableMcpServers {
		// 同步网关stage MCP Server 配置
		result, err = defaultManager.SyncStageMcpConfig(ctx)
		if err != nil {
			log.Fatalf("sync stage mcp servers: %v", err)
		}
		log.Printf("syncing stage mcp servers success, result:%v\n", result)
	}
}
