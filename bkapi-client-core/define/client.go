package define

//go:generate mockgen -source=$GOFILE -destination=./mock/$GOFILE -package=mock BkApiClient,BkApiClientOption

// BkApiClient defines the interface of BkApi client.
type BkApiClient interface {
	// Apply method applies the given options to the client.
	Apply(opts ...BkApiClientOption) error

	// NewOperation method creates a new operation dynamically and apply the given options.
	NewOperation(config OperationConfig, opts ...OperationOption) Operation
}

// BkApiClientOption defines the interface of BkApi client option.
type BkApiClientOption interface {
	// ApplyTo method applies the option to the client.
	ApplyTo(client BkApiClient) error
}
