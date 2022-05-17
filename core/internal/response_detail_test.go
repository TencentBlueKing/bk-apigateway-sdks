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

package internal_test

import (
	"net/http"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResponseDetail", func() {
	It("should create a response detail", func() {
		response := http.Response{
			Header: http.Header{
				"X-Bkapi-Request-Id":    []string{"request-id"},
				"X-Bkapi-Error-Code":    []string{"error-code"},
				"X-Bkapi-Error-Message": []string{"error-message"},
			},
		}
		detail := internal.NewBkApiResponseDetailFromResponse(&response)
		values := detail.Map()

		Expect(values["bkapi_request_id"]).To(Equal(detail.RequestId()))
		Expect(values["bkapi_error_code"]).To(Equal(detail.ErrorCode()))
		Expect(values["bkapi_error_message"]).To(Equal(detail.ErrorMessage()))
		Expect(detail.Cause()).To(Equal(define.ErrBkApiRequest))

		errorMessage := detail.Error()
		Expect(errorMessage).To(ContainSubstring("request-id"))
		Expect(errorMessage).To(ContainSubstring("error-code"))
		Expect(errorMessage).To(ContainSubstring("error-message"))
	})

	It("should return an error", func() {
		detail := internal.NewBkApiResponseDetail("request-id", "error-code", "error-message")
		Expect(detail.GetError()).NotTo(BeNil())
	})

	It("should not return an error", func() {
		detail := internal.NewBkApiResponseDetail("", "", "")
		Expect(detail.GetError()).To(BeNil())
	})
})
