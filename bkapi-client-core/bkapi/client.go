package bkapi

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"github.com/pkg/errors"
	"gopkg.in/h2non/gentleman.v2"
)

// Config is the configuration of BkApi client.
type Config struct {
	Endpoint  string
	AppCode   string
	AppSecret string
	Stage     string
}

// NewBkApiClient creates a new BkApiClient.
func NewBkApiClient(name string, config Config, opts ...define.BkApiClientOption) (*internal.BkApiClient, error) {
	gentlemanClient := gentleman.New().
		URL(config.Endpoint)

	client := internal.NewBkApiClient(name, gentlemanClient, func(name string, request *gentleman.Request) define.Operation {
		return internal.NewOperation(name, request)
	})

	if len(opts) == 0 {
		return client, nil
	}

	err := client.Apply(opts...)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to apply options to client %s", name)
	}

	return client, nil
}
