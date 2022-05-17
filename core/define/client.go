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

package define

//go:generate mockgen -source=$GOFILE -destination=../internal/mock/$GOFILE -package=mock BkApiClient,BkApiClientOption
//go:generate mockgen -destination=../internal/mock/http.go -package=mock net/http RoundTripper
//go:generate mockgen -destination=../internal/mock/io.go -package=mock io ReadCloser

// BkApiClient defines the interface of BkApi client.
type BkApiClient interface {
	// Name method returns the client's name.
	Name() string

	// Apply method applies the given options to the client.
	Apply(opts ...BkApiClientOption) error

	// AddOperationOptions adds the common options to each operation.
	AddOperationOptions(opts ...OperationOption) error

	// NewOperation method creates a new operation dynamically and apply the given options.
	NewOperation(config OperationConfigProvider, opts ...OperationOption) Operation
}

// BkApiClientOption defines the interface of BkApi client option.
type BkApiClientOption interface {
	// ApplyToClient method applies the option to the client.
	ApplyToClient(client BkApiClient) error
}
