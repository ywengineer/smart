package mr_smart

import (
	"github.com/ywengineer/mr.smart/module"
	"go.uber.org/zap"
	"reflect"
)

type Request struct {
	channel     *SocketChannel
	messageCode int
	body        []byte
}

// find  and invoke handler for req
func dispatchRequest(req *Request) {
	hd := module.FindHandlerDefinition(req.messageCode)
	if hd == nil {
		serverLogger.Info("handler definition not found for message code", zap.Int("msgCode", req.messageCode))
		return
	}
	in := hd.CreateIn()
	// decode message
	if err := req.channel.codec.Decode(req.body, in); err != nil {
		// decode failed. close channel
		serverLogger.Info("decode message error. suspicious channel, close it.", zap.Error(err))
		_ = req.channel.Close()
		return
	}
	out := hd.GetMethod().Call([]reflect.Value{reflect.ValueOf(req.channel), reflect.ValueOf(in)})
	// oneway message
	if len(out) == 0 {
		return
	}
	// response
	response := out[0].Interface()
	// execute successful
	if response != nil {
		// send response
		if err := req.channel.Send(out); err != nil {
			serverLogger.Error("send response error", zap.Error(err))
		}
	} else {
		serverLogger.Info("response is nil for handler", zap.Int("msgCode", req.messageCode))
	}
}
