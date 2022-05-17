/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2017 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package bkapi

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
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

// OptSetRequestBody sets the body of the operation.
func OptSetRequestBody(data interface{}) define.OperationOption {
	return NewOperationOption(func(op define.Operation) error {
		op.SetBody(data)
		return nil
	})
}

// OptSetRequestResult sets the result of the operation.
func OptSetRequestResult(result interface{}) define.OperationOption {
	return NewOperationOption(func(op define.Operation) error {
		op.SetResult(result)
		return nil
	})
}

// OptSetRequestPathParams sets the path parameters of the operation.
func OptSetRequestPathParams(params map[string]string) define.OperationOption {
	return NewOperationOption(func(op define.Operation) error {
		op.SetPathParams(params)
		return nil
	})
}
