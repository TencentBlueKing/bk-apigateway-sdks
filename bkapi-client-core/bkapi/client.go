package bkapi

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"github.com/pkg/errors"
	"gopkg.in/h2non/gentleman.v2"
)

func newGentlemanClient(config define.ClientConfig) *gentleman.Client {
	client := gentleman.New().
		URL(config.GetUrl())

	headers := config.GetAuthorizationHeaders()
	if len(headers) > 0 {
		client.SetHeaders(headers)
	}

	return client
}

// NewBkApiClient creates a new BkApiClient.
func NewBkApiClient(apiName string, configProvider define.ClientConfigProvider, opts ...define.BkApiClientOption) (*internal.BkApiClient, error) {
	config := configProvider.Config(apiName)
	gentlemanClient := newGentlemanClient(config)

	client := internal.NewBkApiClient(apiName, gentlemanClient, func(name string, request *gentleman.Request) define.Operation {
		return internal.NewOperation(name, request)
	})

	if len(opts) == 0 {
		return client, nil
	}

	err := client.Apply(opts...)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to apply options to client %s", apiName)
	}

	return client, nil
}
