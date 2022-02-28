package internal

import (
	"gopkg.in/h2non/gentleman.v2/plugin"
)

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

// SimpleOperationOption wrap a operation option that can be used to client
type SimpleOperationOption struct {
	*BkApiClientOption
	*OperationOption
}

// NewSimpleOperationOption creates a new SimpleOperationOption.
func NewSimpleOperationOption(fn func(operation *Operation) error) *SimpleOperationOption {
	opt := &SimpleOperationOption{
		OperationOption: NewOperationOption(fn),
	}
	opt.BkApiClientOption = NewBkApiClientOption(func(client *BkApiClient) error {
		client.operationOptions = append(client.operationOptions, opt)
		return nil
	})

	return opt
}
