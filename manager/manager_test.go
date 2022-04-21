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

package manager_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	apigateway "github.com/TencentBlueKing/bk-apigateway-sdks/apigateway"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	mgr "github.com/TencentBlueKing/bk-apigateway-sdks/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gock.v1"
)

var _ = Describe("Manager", func() {
	var (
		manager *mgr.Manager
		apiName = "testing"
		config  bkapi.ClientConfig
	)

	BeforeEach(func() {
		config = bkapi.ClientConfig{
			Endpoint:  "http://example.com",
			AppCode:   "app_code",
			AppSecret: "app_secret",
		}

		var err error
		manager, err = mgr.NewManager(
			apiName,
			config,
			nil,
			func(configProvider define.ClientConfigProvider, opts ...define.BkApiClientOption) (*apigateway.Client, error) {
				opts = append(opts, bkapi.OptTransport(gock.NewTransport()))
				return apigateway.New(configProvider, opts...)
			},
		)
		Expect(err).To(BeNil())
	})

	It("should load definition from file", func() {
		dir, err := os.MkdirTemp("", "")
		Expect(err).To(BeNil())
		defer os.RemoveAll(dir)

		definitionFile := filepath.Join(dir, "test.yaml")
		Expect(ioutil.WriteFile(
			definitionFile,
			[]byte(`key: {{ data.value }}`),
			0644,
		)).To(Succeed())

		Expect(manager.LoadDefinition(definitionFile, map[string]interface{}{
			"value": "test",
		})).To(Succeed())

		definition := manager.GetDefinition()
		Expect(definition.Get("")).To(Equal(map[string]interface{}{
			"key": "test",
		}))
	})

	It("should return public key", func() {
		gock.New(config.Endpoint).
			Get("/api/v1/apis/testing/public_key/").
			Reply(200).
			JSON(map[string]interface{}{
				"code": 0,
				"data": map[string]interface{}{
					"public_key": "public_key",
				},
			})
		defer gock.Off()

		info, err := manager.GetPublicKey()
		Expect(err).To(BeNil())
		Expect(info).To(Equal(map[string]interface{}{
			"public_key": "public_key",
		}))
	})

	It("should return public key string", func() {
		gock.New(config.Endpoint).
			Get("/api/v1/apis/testing/public_key/").
			Reply(200).
			JSON(map[string]interface{}{
				"code": 0,
				"data": map[string]interface{}{
					"public_key": "public_key",
				},
			})
		defer gock.Off()

		key, err := manager.GetPublicKeyString()
		Expect(err).To(BeNil())
		Expect(key).To(Equal("public_key"))
	})
})
