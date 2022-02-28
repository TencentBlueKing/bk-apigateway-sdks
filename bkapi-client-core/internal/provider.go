package internal

import (
	"bytes"
	"encoding/json"

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
