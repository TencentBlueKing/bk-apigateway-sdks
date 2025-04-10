package gen

import (
	"testing"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_fram/example/router"
)

func TestGenResourceYamlFromSwaggerJson(t *testing.T) {
	GenResourceYamlFromSwaggerJson("../example/docs/swagger.json", router.New())
	// output:
	// paths:
	//  /testapi/get-string-by-int/{some_id}:
	//    get:
	//      consumes:
	//      - application/json
	//      description: get string by ID
	//      operationId: get-string-by-int
	//      parameters:
	//      - description: Some ID
	//        in: path
	//        name: some_id
	//        required: true
	//        type: integer
	//      - description: Some ID
	//        in: body
	//        name: some_id
	//        required: true
	//        schema:
	//          $ref: '#/definitions/api.Pet'
	//      produces:
	//      - application/json
	//      responses:
	//        "200":
	//          description: ok
	//          schema:
	//            type: string
	//      summary: Add a new pet to the store
	//      x-bk-apigateway-resource:
	//        allowApplyPermission: true
	//        authConfig:
	//          appVerifiedRequired: false
	//          resourcePermissionRequired: false
	//          userVerifiedRequired: false
	//        backend:
	//          matchSubpath: false
	//          method: get
	//          path: /testapi/get-string-by-int/{some_id}
	//        enableWebsocket: false
	//        isPublic: true
	//        matchSubpath: false
	//        pluginConfigs:
	//        - type: bk-cors
	//          yaml: |
	//            allow_origins: '*'
	//            allow_methods: '**'
	//            allow_headers: '**'
	//            expose_headers: ''

}
