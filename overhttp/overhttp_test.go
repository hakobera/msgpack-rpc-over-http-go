package overhttp_test

import (
	"testing"

	"github.com/hakobera/msgpack-rpc-over-http-go/overhttp"
)

func TestCall(t *testing.T) {
	url := "http://localhost:9000"
	opts := make(map[string]int32)
	client := overhttp.NewMsgpackClient(url, &opts)

	ret, err := client.Call("add", 1, 2)
	if err != nil {
		t.Errorf("got %v", err)
	}

	if ret != int64(3) {
		t.Errorf("expected 3 but %v", ret)
	}

	t.Log(ret)
}
