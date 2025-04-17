package main

import (
	"net/http"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/example/router"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		petstore.swagger.io
//	@BasePath	/v2

func main() {
	http.ListenAndServe(":8080", router.New())
}
