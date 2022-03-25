package mr_smart

type Request struct {
	request     *socketChannel
	messageCode int
	body        interface{}
}

func getRequestBodyModel(msgCode int) interface{} {
	return nil
}

// find handlers to req and submit task to worker
func dispatchRequest(req *Request) {
	req.request.worker.CtxGo(req.request.ctx, func() {

	})
}
