package define

//go:generate mockgen -source=$GOFILE -destination=./mock/$GOFILE -package=mock BodyProvider

// BodyProvider defines the function to provide the request body.
type BodyProvider interface {
	// ProvideBody method provides the request body, and returns the content length.
	ProvideBody(operation Operation, data interface{}) error
}
