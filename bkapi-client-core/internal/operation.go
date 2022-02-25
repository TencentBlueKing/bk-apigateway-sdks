package internal

import (
	"context"
	"io"
	"net/http"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/pkg/errors"
	"gopkg.in/h2non/gentleman.v2"
	gmctx "gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

// Operation is a wrapper for a request, it allows to set the request options
// and send the request.
type Operation struct {
	name          string
	err           error
	result        interface{}
	resultDecoder func(response *http.Response, result interface{}) error
	request       *gentleman.Request
}

// String returns the operation name.
func (op *Operation) String() string {
	return op.name
}

// Apply method applies the given options to the operation.
func (op *Operation) Apply(opts ...define.OperationOption) define.Operation {
	for _, opt := range opts {
		err := opt.ApplyTo(op)
		if err != nil {
			op.err = errors.WithMessagef(err, "failed to apply option %s", opt)
		}
	}

	return op
}

// SetHeaders used to set the request headers.
func (op *Operation) SetHeaders(headers map[string]string) define.Operation {
	op.request.SetHeaders(headers)

	return op
}

// SetQueryParams used to set the request query parameters.
func (op *Operation) SetQueryParams(params map[string]string) define.Operation {
	op.request.SetQueryParams(params)

	return op
}

// SetPathParams used to set the request path parameters.
func (op *Operation) SetPathParams(params map[string]string) define.Operation {
	op.request.Use(plugin.NewRequestPlugin(func(ctx *gmctx.Context, h gmctx.Handler) {
		ctx.Request.URL.Path = ReplacePlaceHolder(ctx.Request.URL.Path, params)

		h.Next(ctx)
	}))

	return op
}

// SetBody used to set the operation body.
func (op *Operation) SetBody(body io.Reader) define.Operation {
	op.request.Body(body)

	return op
}

// SetBodyProvider used to set the operation body provider.
func (op *Operation) SetBodyProvider(bodyProvider func(operation define.Operation)) define.Operation {
	bodyProvider(op)

	return op
}

// SetResult used to set the operation result.
func (op *Operation) SetResult(result interface{}) define.Operation {
	op.result = result

	return op
}

// SetResultDecoder used to set the operation result decoder.
func (op *Operation) SetResultDecoder(decoder func(response *http.Response, result interface{}) error) define.Operation {
	op.resultDecoder = decoder

	return op
}

// SetContext used to set the request context.
func (op *Operation) SetContext(ctx context.Context) define.Operation {
	op.request.Context.SetCancelContext(ctx)

	return op
}

// Request will send the operation request and return the response.
func (op *Operation) Request() (*http.Response, error) {
	// when the operation already has an error, return it directly
	if op.err != nil {
		return nil, op.err
	}

	response, err := op.request.Send()
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, response.Error
	}

	rawResponse := response.RawResponse
	if op.resultDecoder == nil {
		return rawResponse, nil
	}

	// if the operation has a result decoder, decode the response body
	return rawResponse, op.resultDecoder(rawResponse, op.result)
}

// NewOperation creates a new operation.
func NewOperation(name string, request *gentleman.Request) *Operation {
	return &Operation{
		name:    name,
		request: request,
	}
}

// OperationOption is a wrapper for a operation option.
type OperationOption struct {
	fn func(operation *Operation) error
}

// ApplyTo will check if the operation is valid and apply the option to the operation.
func (o *OperationOption) ApplyTo(op define.Operation) error {
	operation, ok := op.(*Operation)
	if !ok {
		return errors.WithMessagef(
			define.ErrTypeNotMatch, "expected type *Operation, got %T", op,
		)
	}

	return o.fn(operation)
}

// NewOperationOption creates a new operation option.
func NewOperationOption(fn func(operation *Operation) error) define.OperationOption {
	return &OperationOption{
		fn: fn,
	}
}
