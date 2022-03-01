package define

import (
	"context"
	"io"
	"net/http"
)

//go:generate mockgen -source=$GOFILE -destination=../internal/mock/$GOFILE -package=mock Operation,OperationOption

// Operation defines the operation of the API.
type Operation interface {
	// Name method returns the operation's name.
	Name() string

	// Apply method applies the given options to the operation.
	Apply(opts ...OperationOption) Operation

	// SetHeaders method sets the request headers.
	// If the header is already set, it will be overwritten.
	SetHeaders(headers map[string]string) Operation

	// SetQueryParams method sets the request query parameters.
	SetQueryParams(params map[string]string) Operation

	// SetPathParams method sets the request path parameters.
	SetPathParams(params map[string]string) Operation

	// SetBodyReader method sets the request body.
	SetBodyReader(body io.Reader) Operation

	// SetBody method sets the data for body provider
	SetBody(data interface{}) Operation

	// SetBodyProvider method sets the request body provider.
	// A provider not only provides the request body,
	// but also provides the request headers, like Content-Type.
	SetBodyProvider(provider BodyProvider) Operation

	// SetResult method sets the operation result.
	SetResult(result interface{}) Operation

	// SetResultProvider method sets the operation result provider.
	// You can combine multiple decoders into one function,
	// choose the right one by the response status code or content type.
	SetResultProvider(provider ResultProvider) Operation

	// SetContext method sets the request context.
	SetContext(ctx context.Context) Operation

	// SetContentType method sets the request content type.
	SetContentType(contentType string) Operation

	// SetContentEncoding method sets the request content encoding.
	SetContentLength(contentLength int64) Operation

	// Request method sends the operation request and returns the response.
	Request() (*http.Response, error)
}

// OperationOption defines the option of the operation.
type OperationOption interface {
	// ApplyToOperation method applies the option to the operation.
	ApplyToOperation(operation Operation) error
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
