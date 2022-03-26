package mr_smart

import "sync"

type request struct {
	messageCode int
	body        []byte
}

var requestPool = &sync.Pool{
	New: func() interface{} {
		return &request{}
	},
}

func getRequest() *request {
	return requestPool.Get().(*request)
}

func releaseRequest(req *request) {
	req.body = nil
	requestPool.Put(req)
}
