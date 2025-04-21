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

package internal

import (
	"gopkg.in/h2non/gentleman.v2/plugin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
)

// PluginOption wraps a plugin for a client or an operation.
type PluginOption struct {
	*BkApiClientOption
	*OperationOption
}

// ApplyToClient
func (o *PluginOption) ApplyToClient(cli define.BkApiClient) error {
	return o.BkApiClientOption.ApplyToClient(cli)
}

// NewPluginOption creates a new PluginOption.
func NewPluginOption(plugins ...plugin.Plugin) *PluginOption {
	var opt PluginOption
	opt.BkApiClientOption = NewBkApiClientOption(func(cli *BkApiClient) error {
		for _, p := range plugins {
			cli.client.Use(p)
		}

		return nil
	})
	opt.OperationOption = NewOperationOption(func(operation *Operation) error {
		for _, p := range plugins {
			operation.request.Use(p)
		}

		return nil
	})

	return &opt
}
