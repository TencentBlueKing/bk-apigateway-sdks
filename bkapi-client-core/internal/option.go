package internal

import (
	"gopkg.in/h2non/gentleman.v2/plugin"
)

//go:generate mockgen -destination=./mock/plugin.go -package=mock gopkg.in/h2non/gentleman.v2/plugin Plugin

// PluginOption wraps a plugin for a client or an operation.
type PluginOption struct {
	*BkApiClientOption
	*OperationOption
}

// NewPluginOption creates a new PluginOption.
func NewPluginOption(plugins ...plugin.Plugin) *PluginOption {
	var opt PluginOption
	opt.BkApiClientOption = NewBkApiClientOption(func(client *BkApiClient) error {
		client.operationOptions = append(client.operationOptions, &opt)
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
