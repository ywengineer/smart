package smart

import (
	"github.com/ywengineer/smart/codec"
	"github.com/ywengineer/smart/message"
	"github.com/ywengineer/smart/utility"
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

	req, buf := p.Get(), utility.NewLinkBuffer(data)
	_ = json.Decode(buf, req)
	_ = buf.Release()
	t.Logf("%p = %v", req, req)

	p.Put(req)
	req = p.Get()
	//
	data, _ = json.Encode(&message.ForeignMessage{C: 3})
	buf = utility.NewLinkBuffer(data)
	_ = json.Decode(buf, req)
	_ = buf.Release()
	t.Logf("%p = %v", req, req)
}

func TestJSONCodec(t *testing.T) {
	p := sync.Pool{
		New: func() interface{} {
			return &Req{}
		},
	}
	json := codec.Protobuf()

	req, buf := p.Get().(*Req), utility.NewLinkBuffer([]byte(`{"ping": 1, "extra": "abc"}`))
	_ = json.Decode(buf, req)
	_ = buf.Release()
	t.Logf("%p = %v", req, req)

	req.Reset()
	p.Put(req)
	req, buf = p.Get().(*Req), utility.NewLinkBuffer([]byte(`{"ping": 1}`))
	_ = json.Decode(buf, req)
	_ = buf.Release()
	t.Logf("%p = %v", req, req)
}
