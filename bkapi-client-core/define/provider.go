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

import "net/http"

//go:generate mockgen -source=$GOFILE -destination=../internal/mock/$GOFILE -package=mock BodyProvider,ResultProvider

// BodyProvider defines the function to provide the request body.
type BodyProvider interface {
	// ProvideBody method will make the request body by data.
	ProvideBody(operation Operation, data interface{}) error
}

// ResultProvider defines the function to provide the response result.
type ResultProvider interface {
	// ProvideResult method will decode the response body to result.
	ProvideResult(response *http.Response, result interface{}) error
}
