package demo

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
)

// Config is a type alias for bkapi.Config
type Config = bkapi.ClientConfig

// Client extends BkApiClient and defines operations for api gateway demo.
type Client struct {
	define.BkApiClient
}

// Echo call resource echo.
func (c *Client) Echo(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(define.OperationConfig{
		Name:   "echo",
		Method: "POST",
		Path:   "/echo/",
	}, opts...)
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
