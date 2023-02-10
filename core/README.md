# 蓝鲸云 API 客户端(Golang)
蓝鲸 API 网关 SDK 提供了一种基于配置构建 SDK 的方案，本项目基于 gentleman 为该方案提供了底层实现。

## 指引


### 安装
一般情况下，本项目应该作为某个 SDK 的依赖被间接引入，如需直接引用，可以执行：

```golang
go get -u github.com/TencentBlueKing/bk-apigateway-sdks/core@latest
```

### 命名约定：
在 bkapi 包中，作以下命名约定：

| 模式                  | 类型 | 含义                                    | 参考文件  |
| --------------------- | ---- | --------------------------------------- | --------- |
| bkapi.Opt*            | 函数 | 可同时作用于 Client 和 Operation 的选项 | option.go |
| bkapi.*BodyProvider   | 结构 | 提供请求体序列化封装的接口              | json.go   |
| bkapi.*ResultProvider | 结构 | 提供响应体反序列化封装的接口            | json.go   |

### 快速使用
以下基于 [demo](./demo/) 包中的示例：

```golang
import (
	"fmt"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/demo"
)

// 创建客户端，并声明所有结果都使用 Json 格式
client, _ := demo.New(bkapi.ClientConfig{
	Endpoint: "https://httpbin.org/",
}, bkapi.OptJsonResultProvider())

// 创建结果变量
var result demo.AnythingResponse

// 调用接口
_, _ = client.Anything().
	SetResult(&result).
	Request()

// 结果将自动填充到 result 中
fmt.Printf("%#v", result)
```

## 进阶用法
### 启用日志
可通过 `bkapi.ClientConfig` 的 `Logger` 属性来传入日志实现，来捕获相关的流水日志和报错信息，辅助排查问题。
当该属性为空时，默认获取名为 *github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi* 的日志实现。
日志输出会带上请求对应的 Context，可以结合 OTLP 完善链路可观测性。
详见：[github.com/TencentBlueKing/gopkg/logging](https://github.com/TencentBlueKing/gopkg/tree/master/logging)。

### 配置中心
可以使用 `ClientConfigRegistry` 简化每次初始化客户端时都需要传递网关地址，认证信息等重复性工作，`ClientConfigRegistry` 本身已实现成 `ClientConfigProvider`，可直接替代 `ClientConfig` 来使用：

```golang
// 获取默认的配置中心
registry := bkapi.GetGlobalClientConfigRegistry()

// 注册默认的配置（不区分网关）
registry.RegisterDefaultConfig(bkapi.ClientConfig{
	BkApiUrlTmpl: "http://{api_name}.example.com/",
	Stage:        "prod",
})

// 注册指定网关配置
registry.RegisterClientConfig("my-gateway", bkapi.ClientConfig{
	Endpoint:      "http://special-api.example.com/",
	ClientOptions: []define.BkApiClientOption{bkapi.OptJsonResultProvider()},  // 声明这个网关的所有响应都是 JSON
})

// 可直接使用配置中心来初始化客户端
client, _ := New(registry)
```

### Prometheus 指标
*github.com/prometheus/client_golang/prometheus* 模块实现了 Prometheus 插件，启用后可以统计请求过程中的指标：

```golang
// 启用指标插件，可设置指标前缀
prometheus.Enable(prometheus.PrometheusOptions{
	Namespace: "project",
	Subsystem: "module",
})
```

指标一览：
| 名称                            | 类型      | 作用     |
| ------------------------------- | --------- | -------- |
| bkapi_requests_duration_seconds | Histogram | 请求耗时 |
| bkapi_requests_body_bytes       | Histogram | 请求大小 |
| bkapi_responses_body_bytes      | Histogram | 响应大小 |
| bkapi_responses_total           | Counter   | 响应数量 |
| bkapi_failures_total            | Counter   | 失败数量 |

## 定义说明
### 资源封装

Operation 表示一个网关资源封装，方法定义：

| 方法              | 用途                                             |
| ----------------- | ------------------------------------------------ |
| SetHeaders        | 设置请求头                                       |
| SetQueryParams    | 设置请求参数（querystring）                      |
| SetPathParams     | 设置路径变量                                     |
| SetBodyReader     | 设置请求内容（`io.Reader`）                      |
| SetBody           | 设置请求参数（交由请求参数序列化器处理）         |
| SetBodyProvider   | 设置请求参数序列化器（配合 `SetBody`）           |
| SetResult         | 设置响应结构（交由响应序列化器处理）             |
| SetResultProvider | 设置响应序列化器（配合 `SetResult`）             |
| SetContext        | 设置请求上下文                                   |
| SetContentType    | 设置请求 `Content-Type` 头，自定义序列化器可用   |
| SetContentLength  | 设置请求 `Content-Length` 头，自定义序列化器可用 |
| Apply             | 增加额外选项                                     |
| Request           | 发送请求                                         |

### 客户端封装

BkApiClient 表示一个网关封装，方法定义：

| 方法                | 用途                                 |
| ------------------- | ------------------------------------ |
| Apply               | 增加额外选项（作用于 BkApiClient）   |
| AddOperationOptions | 增加资源通用选项（作用于 Operation） |
| NewOperation        | 创建相关的资源封装，并应用通用选项   |

### 客户端配置
客户端配置通过 `bkapi.ClientConfig` 类型来传入，部分参数会自动填充：

| 字段                | 类型                       | 含义           | 必须 | 缺省规则                                                                        |
| ------------------- | -------------------------- | -------------- | ---- | ------------------------------------------------------------------------------- |
| Endpoint            | string                     | 基础地址       | 是   | `"{BkApiUrlTmpl}/{Stage}"`                                                      |
| BkApiUrlTmpl        | string                     | 网关地址模板   | 否   | 环境变量 `BK_API_URL_TMPL`                                                      |
| Stage               | string                     | 环境名称       | 否   | `"prod"`                                                                        |
| AppCode             | string                     | 应用代号       | 否   | 环境变量 `BK_APP_CODE`                                                          |
| AppSecret           | string                     | 应用名称       | 否   | 环境变量 `BK_APP_SECRET`                                                        |
| AccessToken         | string                     | 访问令牌       | 否   |                                                                                 |
| AuthorizationParams | string                     | 额外认证参数   | 否   |                                                                                 |
| Logger              | logging.Logger             | 日志实现       | 否   | `logging.GetLogger("github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi")` |
| ClientOptions       | []define.BkApiClientOption | 通用客户端选项 | 否   |                                                                                 |

注意，Endpoint 和 BkApiUrlTmpl/Stage 选择一种配置方式即可，推荐使用后者。