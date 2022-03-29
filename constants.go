package mr_smart

import (
	"reflect"
	"regexp"
)

const handlerRegexp = "^\\D+([1-9][0-9]*)$"

var TypeSocketChannel = reflect.TypeOf(&SocketChannel{})
var handlerSignatureRegexp = regexp.MustCompile(handlerRegexp)

//var TypeSmartModule = reflect.TypeOf((SmartModule)(nil))
