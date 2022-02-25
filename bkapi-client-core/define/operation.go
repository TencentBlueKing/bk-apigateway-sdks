package define

import (
	"context"
	"io"
	"net/http"
)

//go:generate mockgen -source=$GOFILE -destination=./mock/$GOFILE -package=mock Operation,OperationOption

// Operation defines the operation of the API.
type Operation interface {
	// Apply method applies the given options to the operation.
	Apply(opts ...OperationOption) Operation

	// SetHeaders method sets the request headers.
	// If the header is already set, it will be overwritten.
	SetHeaders(headers map[string]string) Operation

	// SetQueryParams method sets the request query parameters.
	SetQueryParams(params map[string]string) Operation

	// SetPathParams method sets the request path parameters.
	SetPathParams(params map[string]string) Operation

	// SetBody method sets the request body.
	SetBody(body io.Reader) Operation

	// SetBodyProvider method sets the request body provider.
	// A provider not only provides the request body,
	// but also provides the request headers, like Content-Type.
	SetBodyProvider(bodyProvider func(operation Operation)) Operation

	// SetResult method sets the operation result.
	SetResult(result interface{}) Operation

	// SetResultDecoder method sets the operation result decoder.
	// You can combine multiple decoders into one function,
	// choose the right one by the response status code or content type.
	SetResultDecoder(decoder func(response *http.Response, result interface{}) error) Operation

	// SetContext method sets the request context.
	SetContext(ctx context.Context) Operation

	// Request method sends the operation request and returns the response.
	Request() (*http.Response, error)
}

// OperationOption defines the option of the operation.
type OperationOption interface {
	// ApplyTo method applies the option to the operation.
	ApplyTo(Operation) error
}

// OperationConfig used to configure the operation.
type OperationConfig struct {
	// Name is the operation name.
	Name string
	// Method is the HTTP method of the operation.
	Method string
	// Path is the HTTP path of the operation.
	Path string
}
