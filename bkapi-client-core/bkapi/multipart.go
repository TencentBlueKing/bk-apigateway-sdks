package bkapi

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"gopkg.in/h2non/gentleman.v2/plugins/multipart"
)

// MultipartFormFieldsBodyProvider provides request body as multipart form.
type MultipartFormFieldsBodyProvider struct {
	*FunctionalBodyProvider
}

// NewMultipartFormFieldsBodyProvider create a new MultipartFormFieldsBodyProvider
func NewMultipartFormFieldsBodyProvider() *MultipartFormFieldsBodyProvider {
	return &MultipartFormFieldsBodyProvider{
		FunctionalBodyProvider: NewFunctionalBodyProvider(func(operation define.Operation, v interface{}) error {
			values, ok := v.(map[string][]string)
			if !ok {
				return define.ErrorWrapf(define.ErrTypeNotMatch, "expected %T, but got %T", values, v)
			}

			fields := make(map[string]multipart.Values, len(values))
			for k, v := range values {
				fields[k] = multipart.Values(v)
			}

			operation.Apply(internal.NewPluginOption(multipart.Fields(fields)))

			return nil
		}),
	}
}

// MultipartFormBodyProvider provides request body as multipart form.
func MultipartFormBodyProvider() *MultipartFormFieldsBodyProvider {
	return NewMultipartFormFieldsBodyProvider()
}

// OptMultipartFormBodyProvider provides request body as multipart form.
func OptMultipartFormBodyProvider() *MultipartFormFieldsBodyProvider {
	return NewMultipartFormFieldsBodyProvider()
}
