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

package manager

import (
	"github.com/golang-jwt/jwt/v4"

	"github.com/pkg/errors"
)

//
var (
	ErrKidInvalid = errors.New("kid is invalid")
)

// ApigatewayJwtApp represents the request app.
type ApigatewayJwtApp struct {
	AppCode   string `json:"app_code"`
	BkAppCode string `json:"bk_app_code"`
	Verified  bool   `json:"verified"`
}

// ApigatewayJwtUser represents the request user.
type ApigatewayJwtUser struct {
	Username   string `json:"bk_username"`
	SourceType string `json:"source_type"`
	Verified   bool   `json:"verified"`
}

// ApigatewayJwtClaims is the jwt token payload, which carries the information of the request.
type ApigatewayJwtClaims struct {
	jwt.StandardClaims
	ApiName string             `json:"-"`
	App     *ApigatewayJwtApp  `json:"app,omitempty"`
	User    *ApigatewayJwtUser `json:"user,omitempty"`
}

// RsaJwtTokenParser can parse jwt token by rsa algorithm.
type RsaJwtTokenParser struct {
	provider PublicKeyProvider
}

// Parse the jwt token.
func (p *RsaJwtTokenParser) Parse(tokenString string) (ApigatewayJwtClaims, error) {
	var claims ApigatewayJwtClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"]
		if !ok {
			return "", errors.Wrapf(ErrKidInvalid, "kid is not found in jwt header")
		}

		apiName, ok := kid.(string)
		if !ok {
			return "", errors.Wrapf(ErrKidInvalid, "expected kid to be %T but got %T", apiName, kid)
		}

		publicKey, err := p.provider.ProvidePublicKey(apiName)
		if err != nil {
			return "", errors.Wrapf(err, "failed to get public key for %s", apiName)
		}

		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			return pubKey, errors.Wrapf(err, "failed to parse rsa public key for %s", apiName)
		}

		return pubKey, nil
	})

	if err != nil {
		return claims, errors.Wrapf(err, "failed to parse jwt token")
	}

	claims.ApiName = token.Header["kid"].(string)
	return claims, err
}

// NewRsaJwtTokenParser creates a new rsa jwt token parser.
func NewRsaJwtTokenParser(provider PublicKeyProvider) *RsaJwtTokenParser {
	return &RsaJwtTokenParser{
		provider: provider,
	}
}
