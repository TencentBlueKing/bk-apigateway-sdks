package demo

import (
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
)

// Anything : http://httpbin.org/#/Anything/post_anything
func (c *Client) Anything(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "anything",
		Method: "POST",
		Path:   "/anything",
	}, opts...)
}

// StatusCode : http://httpbin.org/#/Status_codes/get_status__codes_
func (c *Client) StatusCode(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "status_code",
		Method: "GET",
		Path:   "/status/{code}",
	}, opts...)
}
