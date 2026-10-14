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

package gen

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/example/router"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/model"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/util"
)

// Under a sub path, the backend paths refer to api_sub_path, which the definition defines.
func TestGenSubPath(t *testing.T) {
	engine := gin.New()
	util.RegisterBkAPIGatewayRoute(engine, "GET", "/users",
		model.NewAPIGatewayResourceConfig(model.ResourceBasicConfig{}),
		func(c *gin.Context) { c.Status(200) })
	swaggerPath := filepath.Join(t.TempDir(), "swagger.json")
	if err := os.WriteFile(swaggerPath, []byte(`{"swagger":"2.0","info":{"title":"users","version":"1"},
		"paths":{"/users":{"get":{"responses":{"200":{"description":"OK"}}}}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		subPath string
		vars    map[string]string
		path    string
	}{
		{"", nil, "/users"},
		{"/my-app", map[string]string{"api_sub_path": "/my-app"}, "{env.api_sub_path}/users"},
	} {
		config := &model.APIConfig{Stage: &model.StageConfig{Name: "prod", BackendSubPath: c.subPath}}
		var definition struct {
			Stages []struct {
				Vars map[string]string `yaml:"vars"`
			} `yaml:"stages"`
		}
		if err := yaml.Unmarshal([]byte(GenDefinitionYaml(config, swaggerPath, engine)), &definition); err != nil {
			t.Fatal(err)
		}
		if vars := definition.Stages[0].Vars; len(vars) != len(c.vars) || vars["api_sub_path"] != c.vars["api_sub_path"] {
			t.Errorf("%q: vars = %v, want %v", c.subPath, vars, c.vars)
		}

		var resources struct {
			Paths map[string]map[string]struct {
				Extension struct {
					Backend model.BackendConfig `yaml:"backend"`
				} `yaml:"x-bk-apigateway-resource"`
			} `yaml:"paths"`
		}
		if err := yaml.Unmarshal([]byte(GenResourceYamlFromSwaggerJson(swaggerPath, engine, c.subPath)), &resources); err != nil {
			t.Fatal(err)
		}
		if path := resources.Paths["/users"]["get"].Extension.Backend.Path; path != c.path {
			t.Errorf("%q: backend path = %q, want %q", c.subPath, path, c.path)
		}
	}
}

func TestGenResourceYamlForMcpWithoutParameters(t *testing.T) {
	engine := gin.New()
	basic := model.ResourceBasicConfig{}
	util.RegisterBkAPIGatewayRoute(engine, "GET", "/mcp-ping",
		model.NewAPIGatewayResourceConfig(basic, basic.WithMcpEnable(true), basic.WithNoneSchema(true)),
		func(c *gin.Context) { c.Status(200) })
	swaggerPath := filepath.Join(t.TempDir(), "swagger.json")
	if err := os.WriteFile(swaggerPath, []byte(`{"swagger":"2.0","info":{"title":"ping","version":"1"},
		"paths":{"/mcp-ping":{"get":{"operationId":"ping","responses":{"200":{"description":"OK"}}}}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	config := &model.APIConfig{Stage: &model.StageConfig{
		Name: "prod", EnableMcpServers: true, McpServerConfigs: []*model.McpServer{{Name: "ping"}},
	}}
	GenDefinitionYaml(config, swaggerPath, engine)
	if !reflect.DeepEqual(config.Stage.McpServerConfigs[0].ResourceNames, []string{"ping"}) {
		t.Fatal("the no-parameter resource was not selected as an MCP tool")
	}
	var resources struct {
		Paths map[string]map[string]struct {
			Extension map[string]any `yaml:"x-bk-apigateway-resource"`
		} `yaml:"paths"`
	}
	if err := yaml.Unmarshal([]byte(GenResourceYamlFromSwaggerJson(swaggerPath, engine, "")), &resources); err != nil {
		t.Fatal(err)
	}
	extension := resources.Paths["/mcp-ping"]["get"].Extension
	if extension["noneSchema"] != true {
		t.Fatalf("expected the gateway's noneSchema field: %v", extension)
	}
}

// The example swagger has the basePath /testapi, under which the routes are served.
func TestGenResourceYamlFromSwaggerJson(t *testing.T) {
	var resources struct {
		Paths map[string]map[string]struct {
			OperationID string                         `yaml:"operationId"`
			Extension   model.APIGatewayResourceConfig `yaml:"x-bk-apigateway-resource"`
		} `yaml:"paths"`
	}
	content := GenResourceYamlFromSwaggerJson("../example/docs/swagger.json", router.New(), "")
	if err := yaml.Unmarshal([]byte(content), &resources); err != nil {
		t.Fatal(err)
	}

	pet := resources.Paths["/pets/{id}/"]["get"]
	if pet.OperationID != "get_pet_by_id" || pet.Extension.IsPublic ||
		pet.Extension.Backend.Path != "/testapi/pets/{id}/" || len(pet.Extension.PluginConfigs) != 2 {
		t.Errorf("unexpected pet resource: %+v", pet)
	}
	product := resources.Paths["/update-product/{product_id}"]["post"]
	if product.OperationID != "update_product_set" || !product.Extension.IsPublic ||
		product.Extension.Backend.Path != "/testapi/update-product/{product_id}" {
		t.Errorf("unexpected product resource: %+v", product)
	}
}
