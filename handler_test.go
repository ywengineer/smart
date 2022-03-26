package mr_smart

import (
	"github.com/ywengineer/mr.smart/codec"
	"github.com/ywengineer/mr.smart/message"
	"sync"
	"testing"
)

func TestHandler(t *testing.T) {
	p := sync.Pool{
		New: func() interface{} {
			return &message.ForeignMessage{}
		},
	}
	json := codec.PBCodec{}

	data, _ := json.Encode(&message.ForeignMessage{C: 1, D: 2})

	req := p.Get()
	json.Decode(data, req)
	t.Logf("%p = %v", req, req)

	p.Put(req)
	req = p.Get()
	data, _ = json.Encode(&message.ForeignMessage{C: 3})
	json.Decode(data, req)
	t.Logf("%p = %v", req, req)
}
