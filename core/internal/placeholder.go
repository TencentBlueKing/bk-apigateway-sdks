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
	"regexp"
	"strings"
)

var placeholderRe *regexp.Regexp

// ReplacePlaceHolder replaces the placeholder with the given string.
func ReplacePlaceHolder(s string, params map[string]string) string {
	return placeholderRe.ReplaceAllStringFunc(
		s, func(placeholder string) string {
			key := strings.Trim(placeholder, "{ }")
			value, ok := params[key]
			if !ok {
				return placeholder
			}

			return value
		},
	)
}

func init() {
	// available placeholder pattern: {param} or { param }
	placeholderRe = regexp.MustCompile(`{\s*.*?\s*}`)
}
