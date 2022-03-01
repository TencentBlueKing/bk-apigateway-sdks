package bkapi

import (
	"encoding/json"
	"io"
)

// JsonMarshalBodyProvider provides request body as json.
type JsonMarshalBodyProvider struct {
	*MarshalBodyProvider
}

// NewJsonMarshalBodyProvider creates a new JsonMarshalBodyProvider with marshal function.
func NewJsonMarshalBodyProvider(marshaler func(v interface{}) ([]byte, error)) *JsonMarshalBodyProvider {
	return &JsonMarshalBodyProvider{
		MarshalBodyProvider: NewMarshalBodyProvider("application/json", marshaler),
	}
}

// JsonBodyProvider creates a new JsonMarshalBodyProvider with default marshal function.
func JsonBodyProvider() *JsonMarshalBodyProvider {
	return NewJsonMarshalBodyProvider(json.Marshal)
}

// OptJsonBodyProvider is a option for json body provider.
func OptJsonBodyProvider() *JsonMarshalBodyProvider {
	return JsonBodyProvider()
}

// JsonUnmarshalResultProvider provides result from json.
type JsonUnmarshalResultProvider struct {
	*UnmarshalResultProvider
}

// NewJsonUnmarshalResultProvider creates a new JsonUnmarshalResultProvider with unmarshal function.
func NewJsonUnmarshalResultProvider(unmarshaler func(body io.Reader, v interface{}) error) *JsonUnmarshalResultProvider {
	return &JsonUnmarshalResultProvider{
		UnmarshalResultProvider: NewUnmarshalResultProvider(unmarshaler),
	}
}

// JsonResultProvider creates a new JsonUnmarshalResultProvider with default unmarshal function.
func JsonResultProvider() *JsonUnmarshalResultProvider {
	return NewJsonUnmarshalResultProvider(func(body io.Reader, v interface{}) error {
		return json.NewDecoder(body).Decode(v)
	})
}

// OptJsonResultProvider is a option for json result provider.
func OptJsonResultProvider() *JsonUnmarshalResultProvider {
	return JsonResultProvider()
}
