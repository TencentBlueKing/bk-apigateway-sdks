package bkapi

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
)

// OperationOption is a wrapper for a operation option.
type OperationOption struct {
	fn func(operation define.Operation) error
}

// ApplyToClient will apply the given options to the client.
func (o *OperationOption) ApplyToClient(client define.BkApiClient) error {
	return client.AddOperationOptions(o)
}

// ApplyToOperation will check if the operation is valid and apply the option to the operation.
func (o *OperationOption) ApplyToOperation(op define.Operation) error {
	return o.fn(op)
}

// NewOperationOption creates a new OperationOption.
func NewOperationOption(fn func(operation define.Operation) error) *OperationOption {
	return &OperationOption{
		fn: fn,
	}
}
