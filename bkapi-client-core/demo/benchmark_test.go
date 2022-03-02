package demo_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/demo"
)

type mockTransport struct {
	response *http.Response
}

func (t mockTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return t.response, nil
}

func newMockTransport() *mockTransport {
	return &mockTransport{
		response: &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Header:     http.Header{},
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"demo": true}`)),
		},
	}
}

func Benchmark_Demo_Request(b *testing.B) {
	client, _ := demo.New(bkapi.ClientConfig{}, bkapi.OptTransport(newMockTransport()))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var result map[string]interface{}
			_, _ = client.Anything().
				SetResultProvider(bkapi.JsonResultProvider()).
				SetResult(&result).
				Request()
		}
	})
}
