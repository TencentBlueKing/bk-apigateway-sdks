# 从 v1 升级到 v2

v2 不兼容 v1，需要 Go 1.26 及以上版本。`manager` 调用网关的 v2 接口，要求网关 1.24.0 及以上版本。

## 模块路径

```bash
go get github.com/TencentBlueKing/bk-apigateway-sdks/v2
go mod tidy
```

| v1 | v2 |
| --- | --- |
| `bk-apigateway-sdks/core/bkapi` | `bk-apigateway-sdks/v2/bkapi`，用法有变化，见下文 |
| `bk-apigateway-sdks/manager` | `bk-apigateway-sdks/v2/manager` |
| `bk-apigateway-sdks/gin_contrib/...` | `bk-apigateway-sdks/v2/gin_contrib/...` |
| `bk-apigateway-sdks/apigateway` | 删除，改用 `manager` 或 `bkapi` |
| `core/define`、`core/prometheus`、`demo` | 删除 |

基于 `core` 生成的网关 SDK 不能只改 import 路径，生成模板和调用代码都要按下文调整。

## bkapi

`bkapi.New` 直接返回 `*gentleman.Client`，请求用 gentleman 的写法构造。`BkApiClient`、`Operation`、各种 provider 和 option、`ClientConfigRegistry` 都已删除。

v1：

```go
client, err := bkapi.NewBkApiClient("demo", bkapi.ClientConfig{
	Endpoint:      "http://special-api.example.com/prod",
	ClientOptions: []define.BkApiClientOption{bkapi.OptJsonResultProvider()},
})
var result QueryUserResponse
_, err = client.NewOperation(
	bkapi.OperationConfig{Name: "query_user", Method: "GET", Path: "/get/{team_id}/user/"},
	bkapi.OptSetRequestPathParams(map[string]string{"team_id": "1"}),
	bkapi.OptSetRequestQueryParam("name", "demo"),
).SetResult(&result).Request()
```

v2：

```go
client, err := bkapi.New(bkapi.Config{Endpoint: "http://special-api.example.com/prod"})
res, err := client.Get().
	AddPath("/get/:team_id/user/").
	Param("team_id", "1").
	SetQuery("name", "demo").
	Use(bkapi.WithContext(ctx)).
	Send()
if err != nil {
	return err
}
var result QueryUserResponse
err = res.JSON(&result)
```

对照：

| v1 | v2 |
| --- | --- |
| 路径参数 `{team_id}` + `OptSetRequestPathParams` | 路径参数 `:team_id` + `Param` |
| `OptSetRequestQueryParam(s)` | `SetQuery`、`SetQueryParams` |
| `OptSetRequestBody` + `OptJsonBodyProvider` | `JSON(body)` |
| `OptSetRequestHeader` | `SetHeader` |
| `SetResult(&result).Request()` | `Send()` 后 `res.JSON(&result)` |
| `SetContext(ctx)` | `Use(bkapi.WithContext(ctx))` |

路径要用 `AddPath` 追加；`Path` 会替换掉 Endpoint 里环境的路径。

### 配置

`ClientConfig` 改为 `bkapi.Config`，`New` 不再读取环境变量：

- 部署在 PaaS 上的应用用 `bkapi.ConfigFromEnv(gatewayName, stageName)` 生成配置，`BK_API_URL_TMPL` 的含义和 v1 相同。不再读取 `BK_API_STAGE_URL_TMPL`。
- 也可以自己填 `Endpoint`，它是网关环境的完整地址，如 `http://bkapi.example.com/api/my-gateway/prod`。
- `AppTenantID` 改名为 `TenantID`。
- `Logger` 改为 `*slog.Logger`，默认 `slog.Default()`。
- 超时和 transport 用 `Config.Timeout`、`Config.Transport` 设置。
- `AuthorizationParams`、`AuthorizationJWT`、`Getenv`、`JsonMarshaler` 删除。

认证方式也有变化：v1 设置了 `AccessToken` 就不再发送 `AppCode` 和 `AppSecret`，v2 会把所有非空的认证字段都放进 `X-Bkapi-Authorization`。

完整的配置说明见 [bkapi README](bkapi/README.md)。

### 错误

非 2xx 响应由 `Send()` 返回 `*bkapi.Error`：

```go
var apiErr *bkapi.Error
if errors.As(err, &apiErr) {
	log.Println(apiErr.StatusCode, apiErr.RequestID, apiErr.Data)
}
```

`Data` 是响应体里的 `error` 字段；没有这个字段时是整个响应体；响应体不是 JSON 时，用 `X-Bkapi-Error-*` 响应头和状态码生成。

连接失败、超时、取消等网络错误原样返回，可以用 `errors.Is(err, context.DeadlineExceeded)` 等方式判断。

`manager.ErrApigatewayRequest` 和 `core/define` 的错误类型已删除。

## manager

构造函数接收 `bkapi.Config`，不再需要 client factory：

```go
mgr, err := manager.NewManagerFrom("my-gateway", manager.ConfigFromEnv(), "definition.yaml")
```

请求网关的方法都把 `ctx` 作为第一个参数，另外有这些变化：

| v1 | v2 |
| --- | --- |
| `Release(version)` | `Release(ctx, version, comment)` |
| `ApplyPermissions(namespace)` | `ApplyPermissions(ctx)`，读取 `apply_permissions` |
| `GrantPermissions()` 返回结果 | `GrantPermissions(ctx)` 只返回 `error` |
| `SyncResourceDocByArchive()` 返回结果 | `SyncResourceDocByArchive(ctx)` 只返回 `error` |
| `SyncPluginConfig(namespace)` | 删除，插件写在 `stages[].plugin_configs` 或资源定义里 |
| `NewDefinitionContext` | 删除，模板先渲染成 YAML 再加载 |

新增 `AddRelatedApps(ctx)` 和 `ResourceVersionExists(ctx, version)`。`SyncStageMcpConfig` 会跳过没有 `mcp_servers` 的环境。

### definition.yaml

`grant_permissions` 推荐使用新字段，旧写法仍然可用：

```yaml
grant_permissions:
  - target_app_code: "app1"   # v1 的 bk_app_code
    grant_dimension: "gateway" # v1 的 api，不填默认 gateway
```

`apply_permissions` 从一个对象改为列表，每一项单独申请，可以用 `gateway_name` 申请其他网关的权限：

```yaml
apply_permissions:
  - gateway_name: "another-gateway"
    grant_dimension: "resource"
    resource_names: ["list_items"]
```

`target_app_code` 默认是当前应用，`applicant` 默认等于 `target_app_code`。

## gin_contrib

- `GatewayConfig.APIType` 改为 `IsOfficial bool`。
- `GatewayConfig` 新增 `Kind`，可选 `normal`（默认）和 `ai`，只在创建网关时生效。
- `McpServer` 新增 `Oauth2PersonalClientEnabled`。
- `GenResourceYamlFromSwaggerJson` 和 `MergeSwaggerConfig` 新增 `subPath` 参数，传应用部署的子路径（如 `/stag--default--app`），资源的后端路径会自动加上这个前缀，不需要的话传 `""`。
- `ApigatewayJwtClaims.ApiName` 改名为 `GatewayName`。
- JWT 中间件在第一次调用 `GatewayJWTAuthMiddleware()` 时读取环境变量，不再在 `init` 时读取。请在加载完 `.env` 之后再创建中间件。

`SyncGinGateway` 的签名不变，失败时仍然 `log.Fatal`。同步流程有几处变化：

- 同步资源时带上 `ResourceDocs.Language`。
- 发布时带上 `Release.Comment`。
- 版本号已存在时，用时间戳替换 build 部分，如 `1.0.0+prod` 变成 `1.0.0+20261012193000`。

使用 bk-apigateway-framework 的 Go 模板创建的项目，按上面的说明修改 import、manager 调用、`IsOfficial` 和 `subPath`。
