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

package manager_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	jwt "github.com/golang-jwt/jwt/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	manager "github.com/TencentBlueKing/bk-apigateway-sdks/v2/manager"
)

var _ = Describe("Jwt", func() {
	var (
		token     string
		provider  *manager.PublicKeySimpleProvider
		parser    *manager.RsaJwtTokenParser
		jwtClaims manager.ApigatewayJwtClaims
	)

	BeforeEach(func() {
		jwtClaims = manager.ApigatewayJwtClaims{
			App: &manager.ApigatewayJwtApp{
				AppCode:  "app_code",
				Verified: true,
			},
			User: &manager.ApigatewayJwtUser{
				Username:   "username",
				SourceType: "default",
				Verified:   true,
			},
		}
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		Expect(err).To(BeNil())
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
		Expect(err).To(BeNil())
		publicKey := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes})

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaims)
		jwtToken.Header["kid"] = "testing"
		token, err = jwtToken.SignedString(key)
		Expect(err).To(BeNil())

		provider = manager.NewPublicKeySimpleProvider(map[string]string{
			"testing": string(publicKey),
		})
		parser = manager.NewRsaJwtTokenParser(provider)
	})

	It("should parse token", func() {
		claims, err := parser.Parse(token)
		Expect(err).To(BeNil())

		Expect(claims.GatewayName).To(Equal("testing"))
		Expect(claims.App).To(Equal(jwtClaims.App))
		Expect(claims.User).To(Equal(jwtClaims.User))
	})
})
