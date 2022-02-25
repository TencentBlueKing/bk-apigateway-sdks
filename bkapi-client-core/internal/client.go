package internal

import (
	"fmt"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

// BkApiClient is a base client for define.
type BkApiClient struct {
	name                string
	client              *gentleman.Client
	commonOperationOpts []define.OperationOption
}

// Name returns the client name.
func (cli *BkApiClient) Name() string {
	return cli.name
}

// NewOperation will create a new operation dynamically and apply the given options.
func (cli *BkApiClient) NewOperation(config define.OperationConfig, opts ...define.OperationOption) define.Operation {
	operation := NewOperation(
		fmt.Sprintf("%s.%s", cli.Name(), config.Name),
		cli.client.Request().Method(config.Method).Path(config.Path),
	)

	operation.Apply(cli.commonOperationOpts...)
	operation.Apply(opts...)

	return operation
}

// Apply method applies the given options to the client.
func (cli *BkApiClient) Apply(opts ...define.BkApiClientOption) {
	for _, opt := range opts {
		opt.ApplyTo(cli)
	}
}

// NewBkApiClient creates a new BkApiClient.
func NewBkApiClient(name string) *BkApiClient {
	return &BkApiClient{
		name:                name,
		client:              gentleman.New(),
		commonOperationOpts: make([]define.OperationOption, 0),
	}
}

// GentlemanClientPluginOption is wrapper for gentleman plugin
type GentlemanClientPluginOption struct {
	plugins []plugin.Plugin
}

// ApplyTo applies the given options to the gentleman client.
func (opt *GentlemanClientPluginOption) ApplyTo(cli define.BkApiClient) {
	client, ok := cli.(*BkApiClient)
	if !ok {
		return
	}

	for _, p := range opt.plugins {
		client.client.Use(p)
	}
}

// NewGentlemanClientPluginOption creates a new gentleman client plugin option.
func NewGentlemanClientPluginOption(plugins ...plugin.Plugin) *GentlemanClientPluginOption {
	return &GentlemanClientPluginOption{
		plugins: plugins,
	}
}
