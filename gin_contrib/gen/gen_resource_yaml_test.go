package gen

import (
	"testing"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/example/router"
)

func TestGenResourceYamlFromSwaggerJson(t *testing.T) {
	GenResourceYamlFromSwaggerJson("../example/docs/swagger.json", router.New())
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
