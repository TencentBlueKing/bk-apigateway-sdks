# manager

蓝鲸 API 网关管理 SDK，提供了基本的注册，同步，发布等功能。

## 功能

- 根据预定义的 YAML 文件进行网关创建，更新，发布及资源同步操作；
- 提供了 JWT token 解析工具，校验接口请求来自于 APIGateway；

## 根据 YAML 同步网关配置
### definition.yaml
用于定义网关资源，为了简化使用，使用以下模型进行处理：

```
+---------------------------------+                +--------------------------------+
|                                 |                |                                |
|                                 |                |  +----------------------+      |
|   ns1:                          |                |  |ns1:                  |      |
|     key: {{data.key}}           |                |  |  key: value_from_data+--+   |             +------------------------------+
|                                 |     Render     |  |                      |  |   |    Load     |                              |
|                                 +--------------->+  +----------------------+  +---------------->+  api(key="value_from_data")  |
|   ns2:                          |                |   ns2:                         |             |                              |
|     key: {{environ.THE_KEY}}    |                |     key: value_from_environ    |             +------------------------------+
|                                 |                |                                |
|                                 |                |                                |
|           Template              |                |              YAML              |
+---------------------------------+                +--------------------------------+
```

definition.yaml 中可以使用 Django 模块语法引用和渲染变量，内置以下变量：
- `environ`：环境变量；
- `data`：命令行自定义变量；

推荐在一个文件中统一进行定义，用命名空间来区分不同资源间的定义：
- `apigateway`：定义网关基本信息；
- `stage`：定义环境信息；
- `plugin_configs`：定义网关插件配置；
- `apply_permissions`：申请网关权限；
- `grant_permissions`：应用主动授权；
- `resource_version`：资源版本信息；
- `release`：定义发布内容；
- `resource_docs`：定义资源文档；

### 使用示例

```go

func SyncGinGateway(baseDir string, apiGatewayName string,
		config *model.APIConfig, delete bool) {
	defaultManager, err := manager.NewManagerFrom(
		apiGatewayName,
		bkapi.ClientConfig{},
		baseDir+"/definition.yaml",
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

	// 生成资源版本
	versionInfo, err := defaultManager.GetLatestResourceVersion()
	if err != nil {
		log.Fatalf("get  gateway resource version: err:%v", err)
		return
	}
	fmt.Printf("gateway resource version:%+v\n", versionInfo)

	newVersion := config.Release.Version

	if len(versionInfo) >= 0 {
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
```

## 解析网关 JWT token
### 选择获取网关公钥方式
解析 JWT token 需要使用网关公钥，内置两种方式：

- `PublicKeySimpleProvider`：直接返回预定义的公钥；
- `PublicKeyMemoryCache`：调用网关接口获取公钥，并缓存一段时间；

此外，可以自行实现 `PublicKeyProvider` 接口，自定义获取网关公钥的方式。

### 解析
选择合适的 `PublicKeyProvider` 实现创建 `RsaJwtTokenParser`：
```golang
jwtParser, err := NewRsaJwtTokenParser(getMyPublicKeyProvider())
claims, err := jwtParser.Parse(jwtToken)
```