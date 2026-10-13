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

package middleware_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/middleware"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/model"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/util"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/manager"
)

// TestGatewayJWTAuthMiddleware checks that the middleware reads the config when it is created, which is after
// the environment is loaded, and fetches the public key only once.
func TestGatewayJWTAuthMiddleware(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	var fetches int
	gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetches++
		fmt.Fprintf(w, `{"data": {"public_key": %q}}`, publicKeyPEM(t, &key.PublicKey))
	}))
	defer gateway.Close()
	t.Setenv("BK_API_URL_TMPL", gateway.URL+"/api/{api_name}/")
	t.Setenv("BK_APP_CODE", "my-app")
	t.Setenv("BK_APP_SECRET", "my-secret")

	gin.SetMode(gin.TestMode)
	router := gin.New()
	// a middleware is created for each route, as in the example router
	for _, path := range []string{"/ping", "/pong"} {
		util.RegisterBkAPIGatewayRoute(router, http.MethodGet, path, model.APIGatewayResourceConfig{},
			middleware.GatewayJWTAuthMiddleware(), func(c *gin.Context) { c.Status(http.StatusNoContent) })
	}

	token := signToken(t, key)
	for _, path := range []string{"/ping", "/pong", "/ping"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set(middleware.BkGatewayJWTHeaderKey, token)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		if res.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d, body: %s", res.Code, http.StatusNoContent, res.Body)
		}
	}
	if fetches != 1 {
		t.Errorf("public key fetched %d times, want 1", fetches)
	}
}

func publicKeyPEM(t *testing.T, key *rsa.PublicKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
}

// signToken signs a token of the gateway "my-gateway" with a verified app.
func signToken(t *testing.T, key *rsa.PrivateKey) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &manager.ApigatewayJwtClaims{
		App: &manager.ApigatewayJwtApp{AppCode: "caller", Verified: true},
	})
	token.Header["kid"] = "my-gateway"
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return signed
}
