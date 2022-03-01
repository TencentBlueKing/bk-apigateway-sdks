package define

//go:generate mockgen -source=$GOFILE -destination=../internal/mock/$GOFILE -package=mock ClientConfig,ClientConfigProvider

// ClientConfig :
type ClientConfig interface {
	GetUrl() string
	GetAuthorizationHeaders() map[string]string
}

type ClientConfigProvider interface {
	Config(apiName string) ClientConfig
}
