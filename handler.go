package mr_smart

import (
	"go.uber.org/zap"
	"reflect"
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
}

func findHandlerDefinition(msgCode int) *handlerDefinition {
	return handlerMap[msgCode]
}

func addHandlerDefinition(def *handlerDefinition) {
	if _, ok := handlerMap[def.messageCode]; ok {
		srvLogger.Warn("handler already exists", zap.Int("msgCode", def.messageCode))
	} else {
		handlerMap[def.messageCode] = def
	}
}

// todo need ObjectPool?
func (hd *handlerDefinition) createIn() interface{} {
	in := hd.inType
	if in.Kind() == reflect.Ptr {
		in = in.Elem()
	}
	return reflect.New(in).Interface()
}

func (hd *handlerDefinition) invoke(channel *SocketChannel, request interface{}) interface{} {
	out := hd.method.Call([]reflect.Value{reflect.ValueOf(channel), reflect.ValueOf(request)})
	if len(out) == 0 {
		return nil
	}
	return out[0].Interface()
}
