package smart

import (
	"github.com/ywengineer/smart/codec"
	"github.com/ywengineer/smart/message"
	"sync"
	"testing"
)

func TestPBCodec(t *testing.T) {
	p := sync.Pool{
		New: func() interface{} {
			return &message.ForeignMessage{}
		},
	}
	json := codec.Protobuf()
	data, _ := json.Encode(&message.ForeignMessage{C: 1, D: 2})

	req := p.Get()
	_ = json.Decode(data, req)
	t.Logf("%p = %v", req, req)

	p.Put(req)
	req = p.Get()
	data, _ = json.Encode(&message.ForeignMessage{C: 3})
	_ = json.Decode(data, req)
	t.Logf("%p = %v", req, req)
}

func TestJSONCodec(t *testing.T) {
	p := sync.Pool{
		New: func() interface{} {
			return &Req{}
		},
	}
	json := codec.Protobuf()

	req := p.Get().(*Req)
	_ = json.Decode([]byte(`{"ping": 1, "extra": "abc"}`), req)
	t.Logf("%p = %v", req, req)

	req.Reset()
	p.Put(req)
	req = p.Get().(*Req)
	_ = json.Decode([]byte(`{"ping": 1}`), req)
	t.Logf("%p = %v", req, req)
}
