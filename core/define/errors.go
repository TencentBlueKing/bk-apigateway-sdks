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

package define

import (
	"errors"

	pkgErrors "github.com/pkg/errors"
)

var (
	// ErrTypeNotMatch defines the error which indicates the type of the value is not match.
	ErrTypeNotMatch = errors.New("type not match")
	// ErrBkApiRequest defines the error which indicates the api request error.
	ErrBkApiRequest = errors.New("bkapi request error")
	// ErrConfigInvalid defines the error which indicates the config is invalid.
	ErrConfigInvalid = errors.New("config invalid")
)

var (
	// ErrorWrapf annotates err with the format specifier and arguments.
	ErrorWrapf = pkgErrors.WithMessagef
	// ErrorCause returns the underlying cause of the error, if possible.
	ErrorCause = pkgErrors.Cause
)

// EnableStackTraceErrorWrapf enables stack trace for ErrorWrapf.
func EnableStackTraceErrorWrapf() {
	SetErrorWrapf(pkgErrors.Wrapf)
}

// SetErrorWrapf sets the ErrorWrapf.
// This function is not thread safe, do not call it in parallel.
func SetErrorWrapf(f func(err error, format string, args ...interface{}) error) {
	ErrorWrapf = f
}

// BkApiRequestError is the error returned by api gateway.
type BkApiRequestError interface {
	error
	RequestId() string
	ErrorCode() string
	ErrorMessage() string
}
