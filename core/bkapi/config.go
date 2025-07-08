/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package bkapi

import (
	"fmt"
	"sync"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
)

// ErrClientConfigRegistryValidationFailed for register validation
var ErrClientConfigRegistryValidationFailed = fmt.Errorf("client config validation failed")

// ClientConfigRegistry manage multiple client configs.
type ClientConfigRegistry struct {
	defaultConfigProvider define.ClientConfigProvider
	configs               sync.Map
}

// ProvideConfig return a client config
func (r *ClientConfigRegistry) ProvideConfig(apiName string) define.ClientConfig {
	provider := r.defaultConfigProvider

	value, ok := r.configs.Load(apiName)
	if ok {
		provider = value.(define.ClientConfigProvider)
	}

	return provider.ProvideConfig(apiName)
}

// RegisterDefaultConfig register default client config
func (r *ClientConfigRegistry) RegisterDefaultConfig(provider define.ClientConfigProvider) error {
	r.defaultConfigProvider = provider
	return nil
}

// RegisterClientConfig register a initialized client config
func (r *ClientConfigRegistry) RegisterClientConfig(apiName string, provider define.ClientConfigProvider) error {
	if apiName == "" {
		return define.ErrorWrapf(ErrClientConfigRegistryValidationFailed, "api name is required")
	}

	r.configs.Store(apiName, provider)
	return nil
}

// NewClientConfigRegistry create a client config registry
func NewClientConfigRegistry() *ClientConfigRegistry {
	var registry ClientConfigRegistry
	err := registry.RegisterDefaultConfig(ClientConfig{})
	if err != nil {
		panic(err)
	}

	return &registry
}

var globalClientConfigRegistry = NewClientConfigRegistry()

// GetGlobalClientConfigRegistry return global client config registry
func GetGlobalClientConfigRegistry() *ClientConfigRegistry {
	return globalClientConfigRegistry
}
