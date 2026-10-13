# bkapi

`bkapi` 是调用蓝鲸 API 网关接口的 Go 客户端。

`bkapi.New` 返回一个 [gentleman](https://github.com/h2non/gentleman) 客户端，已经配置好网关认证、租户、超时、错误处理和日志。请求的构造和响应的读取都使用 gentleman 的 API。

## 从 v1 升级

完整变更见 [v2 迁移指南](../MIGRATION_V2.md)。

## 安装

```bash
go get github.com/TencentBlueKing/bk-apigateway-sdks/v2
```

## 使用

```go
package main

import (
	"context"
	"log"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/bkapi"
)

type User struct {
	Name string `json:"name"`
}

func main() {
	// 1. 创建客户端：调用 my-gateway 网关的 prod 环境
	client, err := bkapi.New(bkapi.ConfigFromEnv("my-gateway", "prod"))
	if err != nil {
		log.Fatal(err)
	}

	// 2. 发送请求：GET {endpoint}/users/admin/
	res, err := client.Get().
		AddPath("/users/:name/").
		Param("name", "admin").
		Use(bkapi.WithContext(context.Background())). // 让请求随 ctx 取消。
		Send()
	if err != nil {
		log.Fatal(err)
	}

	// 3. 解析响应
	var user User
	if err := res.JSON(&user); err != nil {
		log.Fatal(err)
	}
	log.Println(user.Name)
}
```

请求路径使用 `AddPath` 追加到网关地址之后。`Path` 会替换网关地址中的路径，一般不要使用。

更多请求写法：

```go
// 查询参数
client.Get().AddPath("/users/").SetQuery("page", "1").Send()

// JSON 请求体
client.Post().AddPath("/users/").JSON(map[string]any{"name": "admin"}).Send()

// 上传文件
client.Post().AddPath("/docs/").File("file", reader).Send()
```

其他用法参考 [gentleman 文档](https://pkg.go.dev/gopkg.in/h2non/gentleman.v2)。

客户端可以并发使用，通常创建一次后复用即可。

## 响应

响应格式由后端接口决定，网关不会改写，按接口的实际格式读取即可：

```go
// JSON
var user User
err = res.JSON(&user)

// 文本
text := res.String()

// 文件
err = res.SaveToFile("report.csv")
```

例如接口按蓝鲸 API 规范返回 `{"data": ...}` 时，定义对应的结构体：

```go
var body struct {
	Data User `json:"data"`
}
err = res.JSON(&body)
```

## 配置

```go
client, err := bkapi.New(bkapi.Config{
	Endpoint:  "https://bkapi.example.com/api/my-gateway/prod",
	AppCode:   "my-app",
	AppSecret: "my-secret",
})
```

| 字段 | 说明 |
| --- | --- |
| `Endpoint` | 网关环境的地址，必填 |
| `AppCode`、`AppSecret` | 应用 ID 和密钥 |
| `AccessToken` | 用户或应用的 access_token |
| `TenantID` | 租户 ID |
| `Timeout` | 请求超时时间，默认 60 秒 |
| `Transport` | 自定义 `http.RoundTripper` |
| `Logger` | 日志记录器，默认 `slog.Default()` |

部署在蓝鲸 PaaS 上的应用可以使用 `ConfigFromEnv(gatewayName, stageName)` 从环境变量读取配置：

| 字段 | 环境变量 |
| --- | --- |
| `Endpoint` | `BK_API_URL_TMPL`，为空时使用 `BK_API_STAGE_URL_TMPL` |
| `AppCode` | `BK_APP_CODE`、`BKPAAS_APP_ID`、`APP_CODE`，取第一个非空值 |
| `AppSecret` | `BK_APP_SECRET`、`BKPAAS_APP_SECRET`、`SECRET_KEY`，取第一个非空值 |
| `TenantID` | `BKPAAS_APP_TENANT_ID`，为空时使用 `system`；未设置时使用 `BK_APP_TENANT_ID` |

`BK_API_URL_TMPL` 是网关的地址，会在末尾追加环境名；`BK_API_STAGE_URL_TMPL` 是环境的地址，用 `{stage}` 表示环境名。两者的 `{api_name}` 或 `{gateway_name}` 都替换为网关名。例如调用 `ConfigFromEnv("my-gateway", "prod")`：

| 环境变量 | 模板 | 生成地址 |
| --- | --- | --- |
| `BK_API_URL_TMPL` | `https://example.com/api/{api_name}/` | `https://example.com/api/my-gateway/prod` |
| `BK_API_STAGE_URL_TMPL` | `https://{stage}-{api_name}.example.com/` | `https://prod-my-gateway.example.com/` |

## 错误处理

`Send` 返回的错误分为两类：

1. 接口返回非 2xx 状态码时，返回 `*bkapi.Error`：

```go
var apiErr *bkapi.Error
if errors.As(err, &apiErr) {
	log.Println(apiErr.StatusCode) // HTTP 状态码，如 400
	log.Println(apiErr.Data)       // 完整的 error 内容
	log.Println(apiErr.RequestID)  // 请求 ID，可用于在网关中查询请求日志
}
```

新版协议保证存在 `error` 字段，其内容由具体 API 定义。例如：

```json
{"error": {"reason": "参数校验失败", "fields": {"name": "必填"}}}
```

该响应的 `apiErr.Data` 为 `map[string]any`，包含 `reason` 和 `fields`。调用方按具体 API 的约定读取。

没有 `error` 字段的旧格式 JSON 对象响应会完整保留在 `Data` 中。非 JSON 或空响应则使用 `X-Bkapi-Error-Code`、`X-Bkapi-Error-Message` 响应头及 HTTP 状态文本，生成含 `code`、`message` 的兜底对象。

2. 网络错误、超时和取消等错误会原样返回，可以用 `errors.Is` 判断：

```go
if errors.Is(err, context.DeadlineExceeded) {
	// 请求超时
}
```

## 日志

每个响应输出一条日志，包含请求方法、路径、状态码、request_id 和耗时，非 2xx 响应还包含完整错误内容（`error` 字段），请求带有 `traceparent` 头时也会记录。日志不包含请求头、查询参数和请求体。2xx 为 `DEBUG` 级别，4xx 为 `WARN`，其他为 `ERROR`。

网络错误、超时和取消等没有响应的错误只返回给调用方，不输出日志，由调用方结合业务上下文记录。

关闭日志：

```go
bkapi.Config{Logger: slog.New(slog.DiscardHandler)}
```
