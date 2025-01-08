package mr_smart

import (
	"context"
	"github.com/ywengineer/mr.smart/message"
	"github.com/ywengineer/mr.smart/utility"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"reflect"
	"sync"
)

var hManager = &handlerManager{
	_handlerMap: make(map[int]*handlerDefinition, 1000),
}

// handler structure
// code : handler for message code
// name : string
// in(*SocketChannel, request): request must be a ptr
type handlerDefinition struct {
	messageCode int
	name        string
	method      reflect.Value
	inType      reflect.Type // must be ptr
	inPool      *sync.Pool
}

func (hd *handlerDefinition) invoke(channel *SocketChannel, in interface{}) interface{} {
	defer hd.releaseIn(in)
	out := hd.method.Call([]reflect.Value{reflect.ValueOf(channel), reflect.ValueOf(in)})
	if len(out) == 0 {
		return nil
	}
	return out[0].Interface()
}

func (hd *handlerDefinition) releaseIn(in interface{}) {
	// no need to invoke Reset method when in is a protobuf message
	if _, ok := in.(proto.Message); !ok {
		in.(message.Reducible).Reset()
	}
	// release in to object pool
	hd.inPool.Put(in)
}

func (hd *handlerDefinition) getIn() interface{} {
	return hd.inPool.Get()
}

type handlerManager struct {
	// can not use directly
	_handlerMap map[int]*handlerDefinition
}

func (hm *handlerManager) invokeHandler(ctx context.Context, c *SocketChannel, req *request) {
	hd := hm.findHandlerDefinition(req.messageCode)
	if hd == nil {
		utility.DefaultLogger().Info("handler definition not found for message code", zap.Int("msgCode", req.messageCode))
		return
	}
	in := hd.getIn()
	// decode message
	if err := c.codec.Decode(req.body, in); err != nil {
		// decode failed. close channel
		utility.DefaultLogger().Info("decode message error. suspicious channel, close it.", zap.Error(err))
		_ = c.Close()
		hd.releaseIn(in)
	} else if response := hd.invoke(c, in); response != nil {
		if err := c.Send(response); err != nil { // send response
			utility.DefaultLogger().Error("send response error", zap.Error(err))
		}
	} else { // oneway message
		// ignore
	}
}

func (hm *handlerManager) findHandlerDefinition(msgCode int) *handlerDefinition {
	return hm._handlerMap[msgCode]
}

func (hm *handlerManager) addHandlerDefinition(def *handlerDefinition) {
	if _, ok := hm._handlerMap[def.messageCode]; ok {
		utility.DefaultLogger().Warn("handler already exists", zap.Int("msgCode", def.messageCode))
	} else {
		utility.DefaultLogger().Debug("register a new method handler", zap.Int("msgCode", def.messageCode))
		def.inPool = &sync.Pool{
			New: func(hd *handlerDefinition) func() interface{} {
				return func() interface{} {
					in := hd.inType
					if in.Kind() == reflect.Ptr {
						in = in.Elem()
					}
					return reflect.New(in).Interface()
				}
			}(def),
		}
		hm._handlerMap[def.messageCode] = def
	}
}
