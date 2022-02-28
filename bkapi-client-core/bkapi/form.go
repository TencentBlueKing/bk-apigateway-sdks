package bkapi

import (
	"net/url"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"github.com/pkg/errors"
)

// OptFormBodyProvider is a function to set form body
func OptFormBodyProvider() *OperationOption {
	return NewOperationOption(func(operation define.Operation) error {
		operation.SetBodyProvider(internal.NewMarshalBodyProvider(
			"application/x-www-form-urlencoded", func(v interface{}) ([]byte, error) {
				values, ok := v.(map[string][]string)
				if !ok {
					return nil, errors.WithMessagef(define.ErrTypeNotMatch, "expected map[string][]string, but got %T", v)
				}

				forms := url.Values(values)
				return []byte(forms.Encode()), nil
			},
		))

		return nil
	})
}
