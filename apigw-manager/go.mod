module github.com/TencentBlueKing/bk-apigateway-sdks/apigw-manager

go 1.16

replace github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core => ../bkapi-client-core

replace github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-bk-apigateway => ../bkapi-bk-apigateway

require (
	github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-bk-apigateway v0.0.0-00010101000000-000000000000
	github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core v0.0.0-00010101000000-000000000000
	github.com/TencentBlueKing/gopkg v1.0.8
	github.com/flosch/pongo2 v0.0.0-20200913210552-0d938eb266f3
	github.com/flosch/pongo2/v5 v5.0.0
	github.com/golang/mock v1.6.0 // indirect
	github.com/onsi/ginkgo/v2 v2.1.3
	github.com/onsi/gomega v1.18.1
	github.com/pkg/errors v0.9.1
	gopkg.in/h2non/gock.v1 v1.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
