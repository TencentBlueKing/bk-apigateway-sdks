package bkapi

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
)

// Config is the configuration of BkApi client.
type Config struct {
	Endpoint  string
	AppCode   string
	AppSecret string
	Stage     string
}

// NewBkApiClient :
func NewBkApiClient(name string, opts ...define.BkApiClientOption) define.BkApiClient {
	client := internal.NewBkApiClient(name)
	client.Apply(opts...)

	return client
}
