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
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-openapi/spec"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/model"
)

func newTestSwagger() spec.Swagger {
	return spec.Swagger{SwaggerProps: spec.SwaggerProps{Paths: &spec.Paths{Paths: map[string]spec.PathItem{
		"/api/users/list": {PathItemProps: spec.PathItemProps{Get: &spec.Operation{}}},
		"/api/users":      {PathItemProps: spec.PathItemProps{Post: &spec.Operation{}}},
	}}}}
}

func backendOf(t *testing.T, operation *spec.Operation) model.BackendConfig {
	t.Helper()
	config, ok := operation.Extensions["x-bk-apigateway-resource"].(*model.APIGatewayResourceConfig)
	if !ok {
		t.Fatalf("x-bk-apigateway-resource not merged: %v", operation.Extensions)
	}
	return config.Backend
}

func TestMergeSwaggerConfig(t *testing.T) {
	routeMap := map[string]*RouteConfig{
		"/api/users/list:get": {Config: &model.APIGatewayResourceConfig{}},
		"/api/users:post": {Config: &model.APIGatewayResourceConfig{
			Backend: model.BackendConfig{Path: "/api/v2/users"},
		}},
	}

	cases := []struct {
		name     string
		subPath  string
		wantList string
		wantPost string
	}{
		{name: "without sub path", subPath: "", wantList: "/api/users/list", wantPost: "/api/v2/users"},
		{
			name:     "with sub path",
			subPath:  "/stag--default--app/",
			wantList: "/stag--default--app/api/users/list",
			wantPost: "/stag--default--app/api/v2/users",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// merge twice to make sure the sub path is not prefixed repeatedly
			MergeSwaggerConfig(newTestSwagger(), routeMap, c.subPath)
			swagger := MergeSwaggerConfig(newTestSwagger(), routeMap, c.subPath)

			list := backendOf(t, swagger.Paths.Paths["/api/users/list"].Get)
			if list.Path != c.wantList || list.Method != "get" {
				t.Errorf("unexpected list backend: %+v", list)
			}
			post := backendOf(t, swagger.Paths.Paths["/api/users"].Post)
			if post.Path != c.wantPost || post.Method != "post" {
				t.Errorf("unexpected post backend: %+v", post)
			}
		})
	}

	if routeMap["/api/users/list:get"].Config.Backend.Path != "" {
		t.Errorf("the shared route config should not be modified: %+v", routeMap["/api/users/list:get"].Config)
	}
}

func TestMergePlainRouteKeepsSwaggerExtension(t *testing.T) {
	var swagger spec.Swagger
	if err := json.Unmarshal([]byte(`{"paths":{"/health":{"get":{
		"x-bk-apigateway-resource":{"backend":{"method":"get","path":"/custom/health"}}
	}}}}`), &swagger); err != nil {
		t.Fatal(err)
	}
	MergeSwaggerConfig(swagger, map[string]*RouteConfig{"/health:get": {OperationID: "generated"}}, "/app")

	operation := swagger.Paths.Paths["/health"].Get
	want := map[string]any{"backend": map[string]any{"method": "get", "path": "/custom/health"}}
	if got := operation.Extensions["x-bk-apigateway-resource"]; !reflect.DeepEqual(got, want) {
		t.Fatalf("resource = %v, want %v", got, want)
	}
	if operation.ID != "generated" {
		t.Fatalf("operation id = %q", operation.ID)
	}
}
