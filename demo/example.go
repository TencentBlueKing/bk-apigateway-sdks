package demo

import (
	"fmt"
	"log"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
)

func clientExample() {

	// 初始化

	// 获取默认的配置中心
	registry := bkapi.GetGlobalClientConfigRegistry()

	// 注册默认的配置（不区分网关）
	err := registry.RegisterDefaultConfig(bkapi.ClientConfig{
		BkApiUrlTmpl: "http://{api_name}.example.com/",
		Stage:        "prod",
	})
	if err != nil {
		log.Printf("registry default config error: %v", err)
		return
	}

	//// 注册指定网关配置
	//registry.RegisterClientConfig("my-gateway", bkapi.ClientConfig{
	//	Endpoint:      "http://special-api.example.com/",
	//	ClientOptions: []define.BkApiClientOption{bkapi.OptJsonResultProvider()}, // 声明这个网关的所有响应都是 JSON
	//})

	// 可直接使用配置中心来初始化客户端
	// client, _ := New(registry)

	// 创建客户端，并声明所有结果都使用 Json 格式
	client, err := New(bkapi.ClientConfig{
		Endpoint: "https://httpbin.org/",
	}, bkapi.OptJsonResultProvider())

	if err != nil {
		log.Printf("client init error: %v", err)
		return
	}
	// 创建结果变量
	var result AnythingResponse

	// 调用接口(Request()的返回值是：*http.Response,err,看具体情况是否需要处理)

	// 传递路径参数
	_, _ = client.StatusCode(bkapi.OptSetRequestPathParams(map[string]string{
		"code": `200`,
	})).SetResult(&result).Request()

	// 传递query参数
	//_, _ = client.StatusCode(bkapi.OptSetRequestQueryParams(map[string]string{
	//	"code": `200`,
	//})).SetResult(&result).Request()

	// 传递单个query参数
	//_, _ = client.StatusCode(bkapi.OptSetRequestQueryParam("code", `200`)).SetResult(&result).Request()

	// 传递body参数
	_, _ = client.Anything(bkapi.OptSetRequestBody(map[string]string{
		"code": `200`,
	})).SetResult(&result).Request()

	_, _ = client.Anything(bkapi.OptSetRequestBody(
		AnythingRequest{Code: "200"})).SetResult(&result).Request()

	// 传递header参数
	_, _ = client.Anything(
		bkapi.OptSetRequestHeader(
			"X-BKAPI-VERSION", "v3",
		)).SetResult(&result).Request()

	// 结果将自动填充到 result 中
	fmt.Printf("%#v", result)

}
