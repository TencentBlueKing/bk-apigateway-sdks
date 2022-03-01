package define

// BkapiOption combines OperationOption and BkApiClientOption.
type BkapiOption interface {
	OperationOption
	BkApiClientOption
}
