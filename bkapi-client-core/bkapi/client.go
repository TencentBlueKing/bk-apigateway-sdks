package bkapi

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"gopkg.in/h2non/gentleman.v2"
)

// Config is the configuration of BkApi client.
type Config struct {
	Endpoint  string
	AppCode   string
	AppSecret string
	Stage     string
}

// NewBkApiClient :
func NewBkApiClient(name string, config Config, opts ...define.BkApiClientOption) define.BkApiClient {
	gentlemanClient := gentleman.New().
		URL(config.Endpoint)

	client := internal.NewBkApiClient(name, gentlemanClient, func(name string, request *gentleman.Request) define.Operation {
		return internal.NewOperation(name, request)
	})
	client.Apply(opts...)

	return client
}
