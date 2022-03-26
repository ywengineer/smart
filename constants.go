package mr_smart

import (
	"reflect"
)

var TypeSocketChannel = reflect.TypeOf(&SocketChannel{})

//var TypeSmartModule = reflect.TypeOf((SmartModule)(nil))

const HandlerPrefix = "Handler"
const HandlerPrefixLength = len(HandlerPrefix)
