/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TencentBlueKing/gopkg/logging"
	gentleman "gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
	"gopkg.in/h2non/gentleman.v2/plugins/headers"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
)

//go:generate mockgen -destination=../internal/mock/logging.go -package=mock github.com/TencentBlueKing/gopkg/logging Logger

// DefaultUserAgent :
var DefaultUserAgent string

// BkApiClient is a base client for define.
type BkApiClient struct {
	name             string
	logger           logging.Logger
	client           *gentleman.Client
	operationOptions []define.OperationOption
	operationFactory func(name string, client define.BkApiClient, request *gentleman.Request) define.Operation
}

// Name returns the client name.
func (cli *BkApiClient) Name() string {
	return cli.name
}

// Apply method applies the given options to the client.
func (cli *BkApiClient) Apply(opts ...define.BkApiClientOption) error {
	for _, opt := range opts {
		err := opt.ApplyToClient(cli)
		if err != nil {
			return define.ErrorWrapf(
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

func (cli *BkApiClient) logResponse(op define.Operation, response *http.Response) {
	logger := cli.logger
	if logger == nil {
		return
	}

	details := NewBkApiResponseDetailFromResponse(response)
	fields := details.Map()

	ctx := response.Request.Context()
	fields["operation"] = op
	fields["status"] = response.Status
	fields["status_code"] = response.StatusCode

	switch response.StatusCode / 100 {
	case 4:
		logger.WarnContext(ctx, "request error caused by client", fields)
	case 5:
		logger.ErrorContext(ctx, "request error caused by server", fields)
	default:
		logger.DebugContext(ctx, "request success", fields)
	}
}

func (cli *BkApiClient) newGentlemanRequest(config define.OperationConfig) *gentleman.Request {
	return cli.client.Request().
		Method(config.GetMethod()).
		Use(headers.Set("User-Agent", DefaultUserAgent)).
		Use(plugin.NewRequestPlugin(func(c *context.Context, h context.Handler) {
			path := strings.TrimSuffix(c.Request.URL.Path, "/")
			c.Request.URL.Path = fmt.Sprintf("%s/%s", path, strings.TrimPrefix(config.GetPath(), "/"))
			h.Next(c)
		}))
}

func (cli *BkApiClient) newOperationName(config define.OperationConfig) string {
	name := config.GetName()
	if name != "" {
		return name
	}

	return fmt.Sprintf("(%s %s)", config.GetMethod(), config.GetPath())
}

func (cli *BkApiClient) applyOperationOptions(op define.Operation, opts ...define.OperationOption) {
	for _, o := range [][]define.OperationOption{
		cli.operationOptions, opts,
	} {
		if len(o) > 0 {
			op.Apply(o...)
		}
	}
}

// NewOperation will create a new operation dynamically and apply the given options.
func (cli *BkApiClient) NewOperation(
	provider define.OperationConfigProvider,
	opts ...define.OperationOption,
) define.Operation {
	config := provider.ProvideConfig()
	request := cli.newGentlemanRequest(config)
	name := cli.newOperationName(config)
	operation := cli.operationFactory(name, cli, request)

	request.Use(plugin.NewResponsePlugin(func(c *context.Context, h context.Handler) {
		cli.logResponse(operation, c.Response)
		h.Next(c)
	}))

	cli.applyOperationOptions(operation, opts...)

	return operation
}

// NewBkApiClient creates a new BkApiClient.
func NewBkApiClient(
	name string,
	client *gentleman.Client,
	factory func(name string, client define.BkApiClient, request *gentleman.Request) define.Operation,
	config define.ClientConfig,
) (*BkApiClient, error) {
	baseUrl := config.GetUrl()
	if baseUrl == "" {
		return nil, define.ErrorWrapf(define.ErrConfigInvalid, "base url is empty")
	}

	client.URL(baseUrl)

	headers := config.GetAuthorizationHeaders()
	if len(headers) > 0 {
		client.SetHeaders(headers)
	}

	return &BkApiClient{
		name:             name,
		client:           client,
		operationFactory: factory,
		logger:           config.GetLogger(),
		operationOptions: make([]define.OperationOption, 0),
	}, nil
}

// BkApiClientOption is a wrapper for a client option.
type BkApiClientOption struct {
	fn func(client *BkApiClient) error
}

// ApplyToClient will check if the given client is a BkApiClient and apply the option to it.
func (o *BkApiClientOption) ApplyToClient(cli define.BkApiClient) error {
	client, ok := cli.(*BkApiClient)
	if !ok {
		return define.ErrorWrapf(
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

func init() {
	DefaultUserAgent = fmt.Sprintf("%s/%s", define.UserAgent, define.Version)
}
