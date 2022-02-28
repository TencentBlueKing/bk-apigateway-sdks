package bkapi

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
)

// OptJsonBodyProvider is a option for json body provider.
func OptJsonBodyProvider() *OperationOption {
	return NewOperationOption(func(operation define.Operation) error {
		operation.SetBodyProvider(internal.NewJsonBodyProvider())
		return nil
	})
}
