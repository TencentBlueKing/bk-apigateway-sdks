/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2017-2021 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package define

import "github.com/TencentBlueKing/gopkg/logging"

//go:generate mockgen -source=$GOFILE -destination=../internal/mock/$GOFILE -package=mock ClientConfig,ClientConfigProvider,OperationConfig,OperationConfigProvider
//go:generate mockgen -destination=../internal/mock/logging.go -package=mock github.com/TencentBlueKing/gopkg/logging Logger

// ClientConfig used to create a new BkApiClient.
type ClientConfig interface {
	// GetUrl returns the url of the client.
	GetUrl() string
	// GetAuthorizationHeaders returns the authorization headers of the client.
	GetAuthorizationHeaders() map[string]string
	// GetLogger returns the client logger.
	GetLogger() logging.Logger
}

// ClientConfigProvider should provide a ClientConfig instance.
type ClientConfigProvider interface {
	// ProvideConfig returns a ClientConfig instance.
	ProvideConfig(apiName string) ClientConfig
}

// OperationConfig used to configure the operation.
type OperationConfig interface {
	// GetName returns the operation name.
	GetName() string
	// GetMethod returns the HTTP method of the operation.
	GetMethod() string
	// GetPath returns the HTTP path of the operation.
	GetPath() string
}

// OperationConfigProvider should provide a OperationConfig instance.
type OperationConfigProvider interface {
	// ProvideConfig returns a OperationConfig instance.
	ProvideConfig() OperationConfig
}
