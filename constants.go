package smart

import (
	"reflect"
	"regexp"
)

const handlerRegexp = "^\\D+([1-9][0-9]*)$"
const FROM = "from"
const TO = "to"

var TypeSocketChannel = reflect.TypeOf(&SocketChannel{})
var handlerSignatureRegexp = regexp.MustCompile(handlerRegexp)

//var TypeSmartModule = reflect.TypeOf((SmartModule)(nil))
