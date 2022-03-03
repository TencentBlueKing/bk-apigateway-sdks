package apigateway

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
)

// Client : bk-apigateway
type Client struct {
	define.BkApiClient
}

// API Gateway Resources
func New(configProvider define.ClientConfigProvider, opts ...define.BkApiClientOption) (*Client, error) {
	client, err := bkapi.NewBkApiClient("bk-apigateway", configProvider, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		BkApiClient: client,
	}, nil
}
