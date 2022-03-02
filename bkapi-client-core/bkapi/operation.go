package bkapi

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
)

// OperationConfig used to configure the operation.
type OperationConfig struct {
	// Name is the operation name.
	Name string
	// Method is the HTTP method of the operation.
	Method string
	// Path is the HTTP path of the operation.
	Path string
}

// ProvideConfig clone and returns a new OperationConfig.
func (c OperationConfig) ProvideConfig() define.OperationConfig {
	return &c
}

// GetName returns the operation name.
func (c *OperationConfig) GetName() string {
	return c.Name
}

// GetMethod returns the HTTP method of the operation.
func (c *OperationConfig) GetMethod() string {
	return c.Method
}

// GetPath returns the HTTP path of the operation.
func (c *OperationConfig) GetPath() string {
	return c.Path
}

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

// OptSetOperationBody sets the body of the operation.
func OptSetOperationBody(data interface{}) define.OperationOption {
	return NewOperationOption(func(op define.Operation) error {
		op.SetBody(data)
		return nil
	})
}

// OptSetOperationResult sets the result of the operation.
func OptSetOperationResult(result interface{}) define.OperationOption {
	return NewOperationOption(func(op define.Operation) error {
		op.SetResult(result)
		return nil
	})
}
