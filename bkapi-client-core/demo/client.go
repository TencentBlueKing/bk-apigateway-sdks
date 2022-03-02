package demo

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
)

// Client extends BkApiClient and defines operations for http://httpbin.org/.
type Client struct {
	define.BkApiClient
}

// New creates a new Client.
func New(configProvider define.ClientConfigProvider, opts ...define.BkApiClientOption) (*Client, error) {
	client, err := bkapi.NewBkApiClient("demo", configProvider, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		BkApiClient: client,
	}, nil
}
