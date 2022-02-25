package define

//go:generate mockgen -source=$GOFILE -destination=./mock/$GOFILE -package=mock BkApiClient,BkApiClientOption

// BkApiClient defines the interface of BkApi client.
type BkApiClient interface {
	NewOperation(config OperationConfig, opts ...OperationOption) Operation
}

// BkApiClientOption defines the interface of BkApi client option.
type BkApiClientOption interface {
	ApplyTo(BkApiClient)
}
