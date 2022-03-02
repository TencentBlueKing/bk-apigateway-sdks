package define

//go:generate mockgen -source=$GOFILE -destination=../internal/mock/$GOFILE -package=mock BkApiClient,BkApiClientOption
//go:generate mockgen -destination=../internal/mock/http.go -package=mock net/http RoundTripper

// BkApiClient defines the interface of BkApi client.
type BkApiClient interface {
	// Name method returns the client's name.
	Name() string

	// Apply method applies the given options to the client.
	Apply(opts ...BkApiClientOption) error

	// AddOperationOptions adds the common options to each operation.
	AddOperationOptions(opts ...OperationOption) error

	// NewOperation method creates a new operation dynamically and apply the given options.
	NewOperation(config OperationConfigProvider, opts ...OperationOption) Operation
}

// BkApiClientOption defines the interface of BkApi client option.
type BkApiClientOption interface {
	// ApplyToClient method applies the option to the client.
	ApplyToClient(client BkApiClient) error
}
