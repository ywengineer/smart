package smart

import (
	"context"
	"reflect"
	"regexp"
)

const (
	handlerRegexp     = "^\\D+([1-9][0-9]*)$"
	CtxKeySeq         = "sequence"     // value type is int32
	CtxKeyTimestamp   = "timestamp"    // value type is int32
	CtxKeyHeader      = "header"       // value type is map[string][string]
	CtxKeyFrom        = "from"         // value type is string
	CtxKeyService     = "service"      // value type is string, current service name
	CtxKeyFromClient  = "from-client"  // value type is int, connection id
	CtxKeyFromService = "from-service" // value type is string, service name
	CtxKeyToClient    = "to-client"    // value type is int, connection id
	CtxKeyToService   = "to-service"   // value type is string, service name
	CtxKeyTO          = "to"
	HeaderFrom        = CtxKeyFrom // HeaderFrom
)

var TypeSocketChannel = reflect.TypeOf((*Channel)(nil)).Elem()
var handlerSignatureRegexp = regexp.MustCompile(handlerRegexp)
var TypeContext = reflect.TypeOf((*context.Context)(nil)).Elem()

//var TypeSmartModule = reflect.TypeOf((SmartModule)(nil))
