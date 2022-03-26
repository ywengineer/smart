package mr_smart

import (
	"github.com/ywengineer/mr.smart/message"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"reflect"
	"sync"
)

var handlerMap = make(map[int]*handlerDefinition, 1000)

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

func findHandlerDefinition(msgCode int) *handlerDefinition {
	return handlerMap[msgCode]
}

func addHandlerDefinition(def *handlerDefinition) {
	if _, ok := handlerMap[def.messageCode]; ok {
		srvLogger.Warn("handler already exists", zap.Int("msgCode", def.messageCode))
	} else {
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
		handlerMap[def.messageCode] = def
	}
}

func (hd *handlerDefinition) invoke(channel *SocketChannel, in interface{}) interface{} {
	defer func() {
		// no need to invoke Reset method when in is a protobuf message
		if _, ok := in.(proto.Message); !ok {
			in.(message.Reducible).Reset()
		}
		// release in to object pool
		hd.inPool.Put(in)
	}()
	out := hd.method.Call([]reflect.Value{reflect.ValueOf(channel), reflect.ValueOf(in)})
	if len(out) == 0 {
		return nil
	}
	return out[0].Interface()
}
