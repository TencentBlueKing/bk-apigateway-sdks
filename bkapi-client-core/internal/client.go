package internal

import (
	"fmt"
	"strings"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/pkg/errors"
	"gopkg.in/h2non/gentleman.v2"
)

// BkApiClient is a base client for define.
type BkApiClient struct {
	name             string
	client           *gentleman.Client
	operationOptions []define.OperationOption
	operationFactory func(name string, request *gentleman.Request) define.Operation
}

// Name returns the client name.
func (cli *BkApiClient) Name() string {
	return cli.name
}

// Client returns the gentleman client.
func (cli *BkApiClient) Client() *gentleman.Client {
	return cli.client
}

// Apply method applies the given options to the client.
func (cli *BkApiClient) Apply(opts ...define.BkApiClientOption) error {
	for _, opt := range opts {
		err := opt.ApplyToClient(cli)
		if err != nil {
			return errors.WithMessagef(
				err, "failed to apply option %v to client %s", opt, cli.Name(),
			)
		}
	}

	return nil
}

// AddOperationOptions method adds the common options to each operation.
func (cli *BkApiClient) AddOperationOptions(opts ...define.OperationOption) error {
	cli.operationOptions = append(cli.operationOptions, opts...)
	return nil
}

// NewOperation will create a new operation dynamically and apply the given options.
func (cli *BkApiClient) NewOperation(config define.OperationConfig, opts ...define.OperationOption) define.Operation {
	request := cli.client.Request().
		Method(config.Method).
		AddPath(strings.TrimPrefix(config.Path, "/"))

	name := config.Name
	if name == "" {
		name = fmt.Sprintf("(%s %s)", config.Method, config.Path)
	}

	operation := cli.operationFactory(fmt.Sprintf("%s.%s", cli.Name(), name), request)

	for _, o := range [][]define.OperationOption{
		cli.operationOptions, opts,
	} {
		if len(o) > 0 {
			operation.Apply(o...)
		}
	}

	return operation
}

// NewBkApiClient creates a new BkApiClient.
func NewBkApiClient(
	name string,
	client *gentleman.Client,
	factory func(name string, request *gentleman.Request) define.Operation,
) *BkApiClient {
	return &BkApiClient{
		name:             name,
		client:           client,
		operationOptions: make([]define.OperationOption, 0),
		operationFactory: factory,
	}
}

// BkApiClientOption is a wrapper for a client option.
type BkApiClientOption struct {
	fn func(client *BkApiClient) error
}

// ApplyToClient will check if the given client is a BkApiClient and apply the option to it.
func (o *BkApiClientOption) ApplyToClient(cli define.BkApiClient) error {
	client, ok := cli.(*BkApiClient)
	if !ok {
		return errors.WithMessagef(
			define.ErrTypeNotMatch, "expected type %T, got %T", client, cli,
		)
	}

	return o.fn(client)
}

// NewBkApiClientOption creates a new client option.
func NewBkApiClientOption(fn func(client *BkApiClient) error) *BkApiClientOption {
	return &BkApiClientOption{
		fn: fn,
	}
}
