package demo

import (
	"fmt"
	"log"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
)

type QueryUserDemoBodyRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type QueryUserDemoResponse struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Gender  string `json:"gender"`
}

// nolint:unused
func clientExample() {
	// 初始化client

	// 使用默认的全局的注册配置来初始化
	registry := bkapi.GetGlobalClientConfigRegistry()

	// 方式一：注册默认的配置（不区分网关）
	err := registry.RegisterDefaultConfig(bkapi.ClientConfig{
		BkApiUrlTmpl: "http://{api_name}.example.com/", // 网关通用地址
		Stage:        "prod",
	})
	if err != nil {
		log.Printf("registry default config error: %v", err)
		return
	}
	// 使用默认的配置，在创建client的时候再指定网关
	client, err := bkapi.NewBkApiClient("my-gateway", registry)
	if err != nil {
		log.Fatalf("create bkapi client error: %v", err)
		return
	}

	// 方式二：注册指定的网关配置
	err = registry.RegisterClientConfig("my-gateway", bkapi.ClientConfig{
		Endpoint: "http://special-api.example.com/prod", // 具体某个网关地址
		ClientOptions: []define.BkApiClientOption{
			// 设置一些通用的client配置,eg:
			bkapi.OptJsonResultProvider(),                    // 声明这个网关的所有响应都是JSON
			bkapi.OptJsonBodyProvider(),                      // 声明这个网关的body请求都是JSON
			bkapi.OptSetRequestHeader("X-Api-Key", "123456"), // 设置统一的header
		},
	})
	if err != nil {
		log.Fatalf("set bkapi client config error: %v", err)
		return
	}

	client, err = bkapi.NewBkApiClient("my-gateway", registry)
	if err != nil {
		log.Fatalf("create bkapi client error: %v", err)
		return
	}

	// 直接使用自定义的配置来创建
	client, err = bkapi.NewBkApiClient("demo", bkapi.ClientConfig{
		Endpoint: "http://special-api.example.com/prod", // 具体某个网关地址
		ClientOptions: []define.BkApiClientOption{
			// 设置一些通用的client配置,eg:
			bkapi.OptJsonResultProvider(),                    // 声明这个网关的所有响应都是JSON
			bkapi.OptJsonBodyProvider(),                      // 声明这个网关的body请求都是JSON
			bkapi.OptSetRequestHeader("X-Api-Key", "123456"), // 设置统一的header
		},
	})
	if err != nil {
		log.Printf("client init error: %v", err)
		return
	}

	// 创建 api operation
	apiOperation := client.NewOperation(
		// 填充接口配置
		bkapi.OperationConfig{
			Name:   "query_team_user_demo",
			Method: "GET",
			Path:   "/get/{team_id}/user/",
		},
		// 设置header参数
		bkapi.OptSetRequestHeader("X-Bkapi-Header", "demo"),
		// 设置path参数
		bkapi.OptSetRequestPathParams(
			map[string]string{
				"team_id": `1`,
			},
		),
		// 设置query参数
		bkapi.OptSetRequestQueryParam("name", "demo"),
		// 设置body参数: 自定义struct
		bkapi.OptSetRequestBody(QueryUserDemoBodyRequest{Name: "demo"}),
		// 设置body参数： map[string]string
		bkapi.OptSetRequestBody(map[string]string{"name": "demo"}),
	)

	// 创建结果变量
	var result QueryUserDemoResponse

	// 调用接口(Request()的返回值是：*http.Response,err,看具体情况是否需要处理)

	//// 直接通过 api operation传参
	// _,_=apiOperation.SetHeaders(map[string]string{"X-Bkapi-Header": "demo"}).
	//	SetPathParams(map[string]string{"team_id": `1`}).
	//	SetBody(QueryUserDemoBodyRequest{Name: "demo"}).
	//	SetQueryParams(map[string]string{"name": "demo"}).
	//	SetResult(&result).Request()

	_, _ = apiOperation.SetResult(&result).Request()
	// 结果将自动填充到 result 中
	fmt.Printf("%#v", result)
}
