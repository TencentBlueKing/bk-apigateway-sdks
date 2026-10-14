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
		model.NewAPIGatewayResourceConfig(basic, basic.WithMcpEnable(true), basic.WithNonSchema(true)),
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
	if extension["noneSchema"] != true || extension["nonSchema"] != nil {
		t.Fatalf("expected the gateway's noneSchema field: %v", extension)
	}
}

func TestGenResourceYamlFromSwaggerJson(t *testing.T) {
	GenResourceYamlFromSwaggerJson("../example/docs/swagger.json", router.New(), "")
	//`paths:
	//  /testapi/pet/get-pet-by-id/{pet_id}/:
	//    get:
	//      consumes:
	//      - application/json
	//      description: get pet by ID
	//      operationId: get-pet-by-id
	//      parameters:
	//      - description: pet id
	//        in: path
	//        name: pet_id
	//        required: true
	//        type: integer
	//      produces:
	//      - application/json
	//      responses:
	//        "200":
	//          description: OK
	//          schema:
	//            type: string
	//      summary: get a  pet
	//      x-bk-apigateway-resource:
	//        allowApplyPermission: true
	//        authConfig:
	//          appVerifiedRequired: false
	//          resourcePermissionRequired: false
	//          userVerifiedRequired: false
	//        backend:
	//          matchSubpath: false
	//          method: get
	//          path: /testapi/get-pet-by-id/{pet_id}/
	//        enableWebsocket: false
	//        isPublic: false
	//        matchSubpath: false
	//        pluginConfigs:
	//        - type: bk-cors
	//          yaml: |
	//            allow_origins: '*'
	//            allow_methods: '**'
	//            allow_headers: '**'
	//            expose_headers: ''
	//  /testapi/update-product/{product_id}:
	//    post:
	//      consumes:
	//      - application/json
	//      operationId: update-product
	//      parameters:
	//      - description: Product ID
	//        in: path
	//        name: product_id
	//        required: true
	//        type: integer
	//      - description: ' '
	//        in: body
	//        name: productInfo
	//        required: true
	//        schema:
	//          $ref: '#/definitions/api.ProductUpdates'
	//      responses: {}
	//      summary: Update product attributes
	//      x-bk-apigateway-resource:
	//        allowApplyPermission: true
	//        authConfig:
	//          appVerifiedRequired: false
	//          resourcePermissionRequired: false
	//          userVerifiedRequired: false
	//        backend:
	//          matchSubpath: false
	//          method: POST
	//          path: /testapi/product/
	//        enableWebsocket: false
	//        isPublic: true
	//        matchSubpath: false
	//        pluginConfigs:
	//        - type: bk-header-rewrite
	//          yaml: |
	//            set:
	//                X-Test: test
	//            remove:
	//                - X-Test2`
}
