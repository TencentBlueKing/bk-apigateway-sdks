package bkapi

import (
	"bytes"
	"io"
	"net/http"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
)

// MarshalBodyProvider wraps the marshal function to provide the request body.
type MarshalBodyProvider struct {
	contentType string
	marshalFn   func(v interface{}) ([]byte, error)
}

// ContentType returns the Content-Type of the request body.
func (m *MarshalBodyProvider) ContentType() string {
	return m.contentType
}

// ApplyToClient will add to the operation operations.
func (m *MarshalBodyProvider) ApplyToClient(cli define.BkApiClient) error {
	return cli.AddOperationOptions(m)
}

// ApplyToOperation will set the body provider.
func (m *MarshalBodyProvider) ApplyToOperation(op define.Operation) error {
	op.SetBodyProvider(m)
	return nil
}

// ProvideBody method provides the request body, and returns the content length.
func (m *MarshalBodyProvider) ProvideBody(operation define.Operation, data interface{}) error {
	// for most scenarios, a nil data represents an empty body.
	if data == nil {
		return nil
	}

	content, err := m.marshalFn(data)
	if err != nil {
		return define.ErrorWrapf(err, "failed to marshal data to %s", m.contentType)
	}

	operation.
		SetContentType(m.contentType).
		SetContentLength(int64(len(content))).
		SetBodyReader(bytes.NewReader(content))

	return nil
}

// NewMarshalBodyProvider creates a new BodyProvider with the given content type and marshal function.
func NewMarshalBodyProvider(contentType string, marshalFn func(v interface{}) ([]byte, error)) *MarshalBodyProvider {
	return &MarshalBodyProvider{
		contentType: contentType,
		marshalFn:   marshalFn,
	}
}

// UnmarshalResultProvider wraps the unmarshal function to provide result from the response body.
type UnmarshalResultProvider struct {
	unmarshalFn func(body io.Reader, v interface{}) error
}

// ApplyToClient will add to the operation operations.
func (p *UnmarshalResultProvider) ApplyToClient(cli define.BkApiClient) error {
	return cli.AddOperationOptions(p)
}

// ApplyToOperation will set the result provider.
func (p *UnmarshalResultProvider) ApplyToOperation(op define.Operation) error {
	op.SetResultProvider(p)
	return nil
}

// ProvideResult method provides the result from the response body.
func (p *UnmarshalResultProvider) ProvideResult(response *http.Response, result interface{}) error {
	// for most unmarshal functions, a nil receiver is not expected.
	if result == nil {
		return nil
	}

	err := p.unmarshalFn(response.Body, result)
	if err != nil {
		return define.ErrorWrapf(err, "failed to unmarshal response body")
	}

	return nil
}

// NewUnmarshalResultProvider creates a new ResultProvider with the given unmarshal function.
func NewUnmarshalResultProvider(fn func(body io.Reader, v interface{}) error) *UnmarshalResultProvider {
	return &UnmarshalResultProvider{
		unmarshalFn: fn,
	}
}

// FunctionalBodyProvider provides the request body by the given function.
type FunctionalBodyProvider struct {
	fn func(operation define.Operation, data interface{}) error
}

// ApplyToClient will add to the operation operations.
func (p *FunctionalBodyProvider) ApplyToClient(cli define.BkApiClient) error {
	return cli.AddOperationOptions(p)
}

// ApplyToOperation will set the body provider.
func (p *FunctionalBodyProvider) ApplyToOperation(op define.Operation) error {
	op.SetBodyProvider(p)
	return nil
}

// ProvideBody method calls the given function to provide the request body.
func (p *FunctionalBodyProvider) ProvideBody(operation define.Operation, data interface{}) error {
	return p.fn(operation, data)
}

// NewFunctionalBodyProvider creates a new BodyProvider with the given function.
func NewFunctionalBodyProvider(fn func(operation define.Operation, data interface{}) error) *FunctionalBodyProvider {
	return &FunctionalBodyProvider{
		fn: fn,
	}
}
