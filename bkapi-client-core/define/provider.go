package define

import "net/http"

//go:generate mockgen -source=$GOFILE -destination=../internal/mock/$GOFILE -package=mock BodyProvider,ResultProvider

// BodyProvider defines the function to provide the request body.
type BodyProvider interface {
	// ProvideBody method will make the request body by data.
	ProvideBody(operation Operation, data interface{}) error
}

// ResultProvider defines the function to provide the response result.
type ResultProvider interface {
	// ProvideResult method will decode the response body to result.
	ProvideResult(response *http.Response, result interface{}) error
}
