package mr_smart

import (
	"github.com/ywengineer/mr.smart/message"
	"sync"
)

var requestPool = &sync.Pool{
	New: func() interface{} {
		return &message.ProtocolMessage{}
	},
}

func getRequest() *message.ProtocolMessage {
	return requestPool.Get().(*message.ProtocolMessage)
}

func releaseRequest(req *message.ProtocolMessage) {
	req.Reset()
	requestPool.Put(req)
}
