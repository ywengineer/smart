package smart

import (
	"context"
	"reflect"
	"regexp"
)

const (
	handlerRegexp   = "^\\D+([1-9][0-9]*)$"
	CtxKeySeq       = "sequence"  // value type is int32
	CtxKeyTimestamp = "timestamp" // value type is int32
	CtxKeyHeader    = "header"    // value type is map[string][string]
	CtxKeyFrom      = "from"      // value type is string
	CtxKeyTO        = "to"
	//
	HeaderFrom = CtxKeyFrom
)

var TypeSocketChannel = reflect.TypeOf(&SocketChannel{})
var handlerSignatureRegexp = regexp.MustCompile(handlerRegexp)
var TypeContext = reflect.TypeOf((*context.Context)(nil)).Elem()

//var TypeSmartModule = reflect.TypeOf((SmartModule)(nil))
