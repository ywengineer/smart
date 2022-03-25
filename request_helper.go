package mr_smart

type Request struct {
	request     *SocketChannel
	messageCode int
	body        interface{}
}

func getRequestBodyModel(msgCode int) interface{} {
	return nil
}

// find handlers to req and submit task to worker
func dispatchRequest(req *Request) {
	req.request.LaterRun(func() {

	})
}
