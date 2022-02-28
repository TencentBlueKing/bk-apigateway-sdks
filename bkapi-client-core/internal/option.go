package internal

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

// PluginOption wraps a plugin for a client or an operation.
type PluginOption struct {
	*BkApiClientOption
	*OperationOption
}

// ApplyToClient
func (o *PluginOption) ApplyToClient(cli define.BkApiClient) error {
	return o.BkApiClientOption.ApplyToClient(cli)
}

// NewPluginOption creates a new PluginOption.
func NewPluginOption(plugins ...plugin.Plugin) *PluginOption {
	var opt PluginOption
	opt.BkApiClientOption = NewBkApiClientOption(func(cli *BkApiClient) error {
		for _, p := range plugins {
			cli.client.Use(p)
		}

		return nil
	})
	opt.OperationOption = NewOperationOption(func(operation *Operation) error {
		for _, p := range plugins {
			operation.request.Use(p)
		}

		return nil
	})

	return &opt
}
