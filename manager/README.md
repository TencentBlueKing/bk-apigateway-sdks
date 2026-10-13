# manager

`manager` 用于管理蓝鲸 API 网关：根据 `definition.yaml` 同步网关配置、创建并发布资源版本，以及校验网关转发请求中的 JWT。

它通过 [bkapi](../bkapi) 调用 bk-apigateway 的 v2 接口（`/api/v2/sync/`、`/api/v2/open/`）。

## 从 v1 升级

完整变更见 [v2 迁移指南](../MIGRATION_V2.md)。

## 安装

```bash
go get github.com/TencentBlueKing/bk-apigateway-sdks/v2
```

## 同步网关

```go
package main

import (
	"context"
	"log"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/manager"
)

func main() {
	ctx := context.Background()

	mgr, err := manager.NewManagerFrom("my-gateway", manager.ConfigFromEnv(), "definition.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// 同步网关基本信息和环境
	if _, err := mgr.SyncBasicInfo(ctx); err != nil {
		log.Fatal(err)
	}
	if _, err := mgr.SyncStagesConfig(ctx); err != nil {
		log.Fatal(err)
	}

	// 创建资源版本并发布到 definition.yaml 中定义的环境
	if _, err := mgr.CreateResourceVersion(ctx, "1.0.0", "首次发布"); err != nil {
		log.Fatal(err)
	}
	if _, err := mgr.Release(ctx, "1.0.0", "首次发布"); err != nil {
		log.Fatal(err)
	}
}
```

`manager.ConfigFromEnv()` 从环境变量读取应用信息，调用 bk-apigateway 网关的 prod 环境，等同于 `bkapi.ConfigFromEnv("bk-apigateway", "prod")`。需要其他配置时可以直接传入 `bkapi.Config`，字段说明见 [bkapi](../bkapi#配置)。

定义已经在内存中时，使用 `NewManager`：

```go
mgr, err := manager.NewManager("my-gateway", config, manager.NewDefinition(data))
```

完整的同步流程（包括资源、权限、文档、MCP Server）可以参考 gin_contrib 中的 [sync_gin_gateway.go](../gin_contrib/gen/sync_gin_gateway.go)。

## definition.yaml

`definition.yaml` 按命名空间组织网关的定义，每个方法读取对应的命名空间：

| 命名空间 | 说明 | 方法 |
| --- | --- | --- |
| `apigateway` | 网关基本信息 | `SyncBasicInfo` |
| `stages` | 环境列表，包括后端服务、插件配置等 | `SyncStagesConfig`、`Release` |
| `stages[].mcp_servers` | 环境的 MCP Server | `SyncStageMcpConfig` |
| `grant_permissions` | 为其他应用授权访问本网关 | `GrantPermissions` |
| `apply_permissions` | 为本应用申请其他网关的权限 | `ApplyPermissions` |
| `related_apps` | 网关的关联应用 | `AddRelatedApps` |
| `resource_docs` | 资源文档目录 | `SyncResourceDocByArchive` |

示例：

```yaml
spec_version: 2

apigateway:
  description: "示例网关"
  is_public: true
  maintainers:
    - "admin"

stages:
  - name: "prod"
    description: "生产环境"
    backends:
      - name: "default"
        config:
          timeout: 30
          loadbalance: "roundrobin"
          hosts:
            - host: "http://api.example.com"
              weight: 100

grant_permissions:
  - target_app_code: "app1"
    grant_dimension: "gateway"
  - target_app_code: "app2"
    grant_dimension: "resource"
    resource_names: ["get_pet_by_id"]

apply_permissions:
  - gateway_name: "another-gateway"
    grant_dimension: "resource"
    resource_names: ["list_items"]

related_apps:
  - "my-app"

resource_docs:
  basedir: "docs/"
```

权限定义的默认值：

- `grant_dimension` 默认为 `gateway`。
- `gateway_name` 默认为当前网关。
- `apply_permissions` 的 `target_app_code` 默认为当前应用，`applicant` 默认同 `target_app_code`。

资源通过 `SyncResourcesConfig` 单独同步，参数是 `resources.yaml` 的内容，格式见 [示例](../gin_contrib/gen/example/resources.yaml)。

## 校验网关 JWT

网关转发请求时会在 `X-Bkapi-Jwt` 请求头中携带 JWT，后端可以用网关公钥校验它，确认请求来自网关，并获取调用方的应用和用户信息：

```go
provider := manager.NewDefaultPublicKeyMemoryCache(manager.ConfigFromEnv())
parser := manager.NewRsaJwtTokenParser(provider)

claims, err := parser.Parse(r.Header.Get("X-Bkapi-Jwt"))
if err != nil {
	// 校验失败，拒绝请求
}
log.Println(claims.GatewayName, claims.App.AppCode, claims.User.Username)
```

`claims.App`、`claims.User` 在 JWT 不包含对应信息时为 `nil`，使用前需要判断。

公钥的获取方式：

- `NewDefaultPublicKeyMemoryCache`：调用网关接口获取公钥，并缓存 12 小时。
- `NewPublicKeySimpleProvider`：使用预先配置的公钥，key 为网关名。
- 实现 `PublicKeyProvider` 接口，自定义获取方式。

使用 gin 时可以直接使用 gin_contrib 中的 [GatewayJWTAuthMiddleware](../gin_contrib/middleware/apigateway_jwt.go)。它在第一次创建时读取环境变量，请在环境变量加载之后再创建。多次创建共用同一份公钥缓存。

## 错误处理

调用网关接口失败时返回 `*bkapi.Error`，包含状态码、服务端错误内容和请求 ID，详见 [bkapi](../bkapi#错误处理)。
