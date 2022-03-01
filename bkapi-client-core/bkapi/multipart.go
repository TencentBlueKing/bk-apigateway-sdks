package bkapi

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"github.com/pkg/errors"
	"gopkg.in/h2non/gentleman.v2/plugins/multipart"
)

// OptMultipartFormBodyProvider :
func OptMultipartFormBodyProvider() *internal.OperationOption {
	return internal.NewOperationOption(func(operation *internal.Operation) error {
		request := operation.GetGentlemanRequest()

		operation.SetBodyProvider(internal.NewFunctionalBodyProvider(func(op define.Operation, v interface{}) error {
			values, ok := v.(map[string][]string)
			if !ok {
				return errors.WithMessagef(define.ErrTypeNotMatch, "expected map[string][]string, but got %T", v)
			}

			fields := make(map[string]multipart.Values, len(values))
			for k, v := range values {
				fields[k] = multipart.Values(v)
			}

			request.Use(multipart.Fields(fields))

			return nil
		}))

		return nil
	})
}
