package smart

import (
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"gitee.com/ywengineer/smart/codec"
	"gitee.com/ywengineer/smart/message"
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

	req, buf := p.Get(), utilk.NewLinkBuffer(data)
	_ = json.Decode(buf, req)
	_ = buf.Release()
	t.Logf("%p = %v", req, req)

	p.Put(req)
	req = p.Get()
	//
	data, _ = json.Encode(&message.ForeignMessage{C: 3})
	buf = utilk.NewLinkBuffer(data)
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
	json := codec.Json()

	req, buf := p.Get().(*Req), utilk.NewLinkBuffer([]byte(`{"ping": 1, "extra": "abc"}`))
	_ = json.Decode(buf, req)
	_ = buf.Release()
	t.Logf("%p = %v", req, req)

	req.Reset()
	p.Put(req)
	req, buf = p.Get().(*Req), utilk.NewLinkBuffer([]byte(`{"ping": 1}`))
	_ = json.Decode(buf, req)
	_ = buf.Release()
	t.Logf("%p = %v", req, req)
}
