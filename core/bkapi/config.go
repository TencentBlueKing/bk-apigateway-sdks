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
	"fmt"
	"sync"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
)

//
var (
	// ErrClientConfigRegistryValidationFailed for register validation
	ErrClientConfigRegistryValidationFailed = fmt.Errorf("client config validation failed")
)

const clientConfigRegistryDefaultKey = ""

// ClientConfigRegistry manage multiple client configs.
type ClientConfigRegistry struct {
	configs sync.Map
}

func (r *ClientConfigRegistry) getConfig(apiName string) define.ClientConfig {
	for _, key := range []string{apiName, clientConfigRegistryDefaultKey} {
		config, ok := r.configs.Load(key)
		if ok {
			return config.(define.ClientConfig)
		}
	}

	return nil
}

// ProvideConfig return a client config
func (r *ClientConfigRegistry) ProvideConfig(apiName string) define.ClientConfig {
	config := r.getConfig(apiName)

	switch realConfig := config.(type) {
	case *ClientConfig:
		// try to copy a new one
		newConfig := *realConfig
		// for default config
		newConfig.setApiName(apiName)
		return &newConfig
	default:
		return config
	}
}

func (r *ClientConfigRegistry) registerClientConfig(apiName string, provider define.ClientConfigProvider) error {
	config := provider.ProvideConfig(apiName)
	r.configs.Store(apiName, config)
	return nil
}

// RegisterDefaultConfig register default client config
func (r *ClientConfigRegistry) RegisterDefaultConfig(provider define.ClientConfigProvider) error {
	return r.registerClientConfig(clientConfigRegistryDefaultKey, provider)
}

// RegisterClientConfig register a initialized client config
func (r *ClientConfigRegistry) RegisterClientConfig(apiName string, provider define.ClientConfigProvider) error {
	if apiName == "" {
		return define.ErrorWrapf(ErrClientConfigRegistryValidationFailed, "api name is required")
	}

	return r.registerClientConfig(apiName, provider)
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
