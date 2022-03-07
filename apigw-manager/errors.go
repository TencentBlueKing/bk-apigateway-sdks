package manager

import "errors"

//
var (
	ErrNotFound          = errors.New("not found")
	ErrApigatewayRequest = errors.New("apigateway request error")
)
