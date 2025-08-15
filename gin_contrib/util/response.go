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

package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BadRequestError ...
const (
	BadRequestError   = "BadRequest"
	UnauthorizedError = "Unauthorized"
	ForbiddenError    = "Forbidden"
	NotFoundError     = "NotFound"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	System  string `json:"system"` // 错误来源系统
}

// ErrorResponse  ...
type ErrorResponse struct {
	Error Error `json:"error"`
}

// BadRequestErrorJSONResponse ...
var (
	UnauthorizedJSONResponse = NewErrorJSONResponse(UnauthorizedError, http.StatusUnauthorized)
)

func NewErrorJSONResponse(
	errorCode string,
	statusCode int,
) func(c *gin.Context, message string) {
	return func(c *gin.Context, message string) {
		c.JSON(statusCode, ErrorResponse{Error: Error{
			Code:    errorCode,
			Message: message,
			System:  "bk-apigateway-sdks",
		}})
	}
}
