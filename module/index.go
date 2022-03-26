package module

import (
	"reflect"
)

// handler structure
// code : handler for message code
// name : string
// in(*SocketChannel, request): request must be a ptr
// out(response, error) : response must be a ptr
type HandlerDefinition struct {
	messageCode int
	name        string
	method      reflect.Value
	inType      reflect.Type // must be ptr
	outType     reflect.Type // must be ptr or nil
}

func FindHandlerDefinition(msgCode int) *HandlerDefinition {
	return nil
}

func (hd *HandlerDefinition) NewIn() interface{} {
	in := hd.inType
	if in.Kind() == reflect.Ptr {
		in = in.Elem()
	}
	return reflect.New(in).Interface()
}

func (hd *HandlerDefinition) GetMethod() reflect.Value {
	return hd.method
}
