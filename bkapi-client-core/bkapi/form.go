package bkapi

import (
	"net/url"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/pkg/errors"
)

// FormMarshalBodyProvider provides request body as urlencoded form.
type FormMarshalBodyProvider struct {
	*MarshalBodyProvider
}

// NewFormMarshalBodyProvider creates a new FormMarshalBodyProvider with marshal function.
func NewFormMarshalBodyProvider(marshaler func(v interface{}) ([]byte, error)) *FormMarshalBodyProvider {
	return &FormMarshalBodyProvider{
		MarshalBodyProvider: NewMarshalBodyProvider("application/x-www-form-urlencoded", marshaler),
	}
}

// FormBodyProvider is a function to set form body from map[string][]string
func FormBodyProvider() *FormMarshalBodyProvider {
	return NewFormMarshalBodyProvider(func(v interface{}) ([]byte, error) {
		values, ok := v.(map[string][]string)
		if !ok {
			return nil, errors.WithMessagef(define.ErrTypeNotMatch, "expected %T, but got %T", values, v)
		}

		forms := url.Values(values)
		return []byte(forms.Encode()), nil
	})
}

// OptFormBodyProvider is a function to set form body
func OptFormBodyProvider() *FormMarshalBodyProvider {
	return FormBodyProvider()
}
