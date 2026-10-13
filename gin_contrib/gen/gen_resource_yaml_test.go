package gen

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/example/router"
)

func TestGenResourceYamlForPlainGinRoute(t *testing.T) {
	engine := gin.New()
	engine.GET("/health", func(c *gin.Context) { c.Status(200) })
	swaggerPath := filepath.Join(t.TempDir(), "swagger.json")
	if err := os.WriteFile(swaggerPath, []byte(`{"swagger":"2.0","info":{"title":"health","version":"1"},
		"paths":{"/health":{"get":{"responses":{"200":{"description":"OK"}}}}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, subPath := range []string{"", "/stag--default--app", "/stag--default--app/"} {
		t.Run(subPath, func(t *testing.T) {
			content := GenResourceYamlFromSwaggerJson(swaggerPath, engine, subPath)
			var definition struct {
				Paths map[string]map[string]struct {
					OperationID string         `yaml:"operationId"`
					Extension   map[string]any `yaml:"x-bk-apigateway-resource"`
				} `yaml:"paths"`
			}
			if err := yaml.Unmarshal([]byte(content), &definition); err != nil {
				t.Fatal(err)
			}
			operation := definition.Paths["/health"]["get"]
			if operation.OperationID == "" {
				t.Fatal("missing operation ID")
			}
			if subPath == "" {
				if operation.Extension != nil {
					t.Fatalf("unexpected defaults: %v", operation.Extension)
				}
				return
			}
			want := map[string]any{"backend": map[string]any{
				"method": "get", "path": "/stag--default--app/health",
			}}
			if !reflect.DeepEqual(operation.Extension, want) {
				t.Fatalf("extension = %v, want %v", operation.Extension, want)
			}
		})
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
