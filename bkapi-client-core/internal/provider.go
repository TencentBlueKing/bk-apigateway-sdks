package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/pkg/errors"
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

// ProvideBody method provides the request body, and returns the content length.
func (m *MarshalBodyProvider) ProvideBody(operation define.Operation, data interface{}) error {
	content, err := m.marshalFn(data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal data to %s", m.contentType)
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

// NewJsonBodyProvider creates a json BodyProvider.
func NewJsonBodyProvider() *MarshalBodyProvider {
	return NewMarshalBodyProvider("application/json", json.Marshal)
}

// UnmarshalResultProvider wraps the unmarshal function to provide result from the response body.
type UnmarshalResultProvider struct {
	unmarshalFn func(body io.Reader, v interface{}) error
}

// ProvideResult method provides the result from the response body.
func (p *UnmarshalResultProvider) ProvideResult(response *http.Response, result interface{}) error {
	err := p.unmarshalFn(response.Body, result)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal response body")
	}

	return nil
}

// NewUnmarshalResultProvider creates a new ResultProvider with the given unmarshal function.
func NewUnmarshalResultProvider(fn func(body io.Reader, v interface{}) error) *UnmarshalResultProvider {
	return &UnmarshalResultProvider{
		unmarshalFn: fn,
	}
}

// NewJsonResultProvider creates a json ResultProvider.
func NewJsonResultProvider() *UnmarshalResultProvider {
	return NewUnmarshalResultProvider(func(body io.Reader, v interface{}) error {
		return json.NewDecoder(body).Decode(v)
	})
}

// FunctionalBodyProvider :
type FunctionalBodyProvider struct {
	fn func(operation define.Operation, data interface{}) error
}

// ProvideBody :
func (p *FunctionalBodyProvider) ProvideBody(operation define.Operation, data interface{}) error {
	return p.fn(operation, data)
}

// NewFunctionalBodyProvider :
func NewFunctionalBodyProvider(fn func(operation define.Operation, data interface{}) error) *FunctionalBodyProvider {
	return &FunctionalBodyProvider{
		fn: fn,
	}
}
