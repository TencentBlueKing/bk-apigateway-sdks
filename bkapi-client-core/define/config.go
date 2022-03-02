package define

//go:generate mockgen -source=$GOFILE -destination=../internal/mock/$GOFILE -package=mock ClientConfig,ClientConfigProvider,OperationConfig,OperationConfigProvider

// ClientConfig used to create a new BkApiClient.
type ClientConfig interface {
	// GetUrl returns the url of the client.
	GetUrl() string
	// GetAuthorizationHeaders returns the authorization headers of the client.
	GetAuthorizationHeaders() map[string]string
}

// ClientConfigProvider should provide a ClientConfig instance.
type ClientConfigProvider interface {
	// ProvideConfig returns a ClientConfig instance.
	ProvideConfig(apiName string) ClientConfig
}

// OperationConfig used to configure the operation.
type OperationConfig interface {
	// GetName is the operation name.
	GetName() string
	// GetMethod is the HTTP method of the operation.
	GetMethod() string
	// GetPath is the HTTP path of the operation.
	GetPath() string
}

// OperationConfigProvider should provide a OperationConfig instance.
type OperationConfigProvider interface {
	// ProvideConfig returns a OperationConfig instance.
	ProvideConfig() OperationConfig
}
