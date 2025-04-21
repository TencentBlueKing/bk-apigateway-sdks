package gen

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/model"
	"github.com/TencentBlueKing/bk-apigateway-sdks/manager"
)

func SyncGinGateway(baseDir string, apiGatewayName string,
		config *model.APIConfig, delete bool,
) {
	defaultManager, err := manager.NewManagerFrom(
		apiGatewayName,
		bkapi.ClientConfig{},
		strings.TrimSuffix(baseDir, "/")+"/definition.yaml",
	)
	if err != nil {
		log.Fatal("Error creating default manager:", err)
		return
	}

	// 同步网关基础信息
	info, err := defaultManager.SyncBasicInfo()
	if err != nil {
		log.Fatalf("syncing gateway basic info: err:%v", err)
		return
	}
	log.Printf("syncing gateway basic info success, info:%v\n", info)

	// 同步网关环境信息
	result, err := defaultManager.SyncStagesConfig()
	if err != nil {
		log.Fatalf("syncing gateway stage config: err:%v", err)
		return
	}
	log.Printf("syncing gateway stage config success, result:%v\n", result)

	// 同步网关资源信息
	resourceFile, err := os.ReadFile(baseDir + "/resources.yaml")
	if err != nil {
		log.Fatal("Error reading resources file:", err)
		return
	}
	log.Printf("call sync_apigw_resources with resources:%s\n", resourceFile)

	result, err = defaultManager.SyncResourcesConfig(map[string]interface{}{
		"content":  string(resourceFile),
		"delete":   delete,
		"language": config.ResourceDocs.Language,
	})
	if err != nil {
		log.Fatalf("syncing gateway resource config: err:%v", err)
		return
	}
	log.Printf("syncing gateway resource config success, result:%v\n", result)

	// 同步授权信息
	result, err = defaultManager.GrantPermissions()
	if err != nil {
		log.Fatalf("syncing gateway resource config: err:%v", err)
		return
	}
	log.Printf("syncing gateway resource config success, result:%v\n", result)

	// 同步资源文档信息
	if config.ResourceDocs.BaseDir != "" {
		result, err = defaultManager.SyncResourceDocByArchive()
		if err != nil {
			log.Fatalf("syncing gateway resource doc: err:%v", err)
			return
		}
		log.Printf("syncing gateway resource doc success, result:%v\n", result)
	}

	// 生成资源版本
	versionInfo, err := defaultManager.GetLatestResourceVersion()
	if err != nil {
		log.Fatalf("get  gateway resource version: err:%v", err)
		return
	}
	fmt.Printf("gateway resource version:%+v\n", versionInfo)

	newVersion := config.Release.Version

	if len(versionInfo) > 0 {
		oldVersion := versionInfo["version"].(string)
		if strings.Contains(oldVersion, newVersion) {
			newVersion = fmt.Sprintf("%s+%s", newVersion, time.Now().Format("20060102150405"))
		}
	}
	result, err = defaultManager.CreateResourceVersion(newVersion, config.Release.Comment)
	if err != nil {
		log.Fatalf("create gateway resource version: err:%v", err)
		return
	}
	log.Printf("create gateway resource version success, result:%v\n", result)
	// 发布资源版本
	if !config.Release.NoPub {
		result, err = defaultManager.Release(newVersion)
		if err != nil {
			log.Fatalf("release gateway resource version: err:%v", err)
			return
		}
		log.Printf("release gateway resource version success, result:%v\n", result)
	}
}
