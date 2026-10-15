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

package main

import (
	"net/http"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/example/router"
)

//	@title			gin_contrib Example API
//	@version		1.0
//	@description	An example of registering gin routes as the resources of BlueKing API Gateway.

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@BasePath	/testapi

func main() {
	_ = http.ListenAndServe(":8080", router.New())
}
