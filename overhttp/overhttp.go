package overhttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ugorji/go/codec"
)

const (
	REQUEST  int64 = 0 // [0, msgid, method, params]
	RESPONSE int64 = 1 // [1, msgid, error, result]
	NOTIFY   int64 = 2 // [2, method, param]
)

const (
	NO_METHOD_ERROR int64 = 0x01
	ARGUMENT_ERROR  int64 = 0x02
)

type MsgpackRPCOverHTTPClient struct {
	Url               string
	SeqId             uint32
	ConnectionTimeout int32
	SendTimeout       int32
	ReceiveTimeout    int32
}

func NewMsgpackClient(url string, options *map[string]int32) *MsgpackRPCOverHTTPClient {
	client := new(MsgpackRPCOverHTTPClient)
	client.Url = url
	client.SeqId = 0
	return client
}

func (c *MsgpackRPCOverHTTPClient) Call(method string, args ...interface{}) (interface{}, error) {
	return c.SendRequest(method, args...)
}

func (c *MsgpackRPCOverHTTPClient) SendRequest(method string, param ...interface{}) (interface{}, error) {
	data, err := c.CreateRequestBody(method, param...)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(*data)
	resp, err := http.Post(c.Url, "application/x-msgpack", reader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return c.GetResult(&body)
}

func (c *MsgpackRPCOverHTTPClient) CreateRequestBody(method string, param ...interface{}) (*[]byte, error) {
	msgid := c.SeqId
	c.SeqId++
	if c.SeqId >= (1 << 31) {
		c.SeqId = 0
	}

	v := make([]interface{}, 4)
	v[0] = REQUEST
	v[1] = msgid
	v[2] = method
	v[3] = param[:]

	var (
		b  []byte
		mh codec.MsgpackHandle
	)

	enc := codec.NewEncoderBytes(&b, &mh)
	err := enc.Encode(v)

	return &b, err
}

func (c *MsgpackRPCOverHTTPClient) GetResult(body *[]byte) (interface{}, error) {
	var (
		v  []interface{}
		mh codec.MsgpackHandle
	)

	dec := codec.NewDecoderBytes(*body, &mh)
	err := dec.Decode(&v)
	if err != nil {
		return nil, err
	}

	if v[0] != RESPONSE {
		return nil, fmt.Errorf("Unknown message type %s", v[0])
	}

	if v[2] != nil {
		return nil, fmt.Errorf("%s", v[2])
	}

	return v[3], nil
}
