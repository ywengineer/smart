package smart

import (
	"context"
	"github.com/go-spring/spring-core/gs"
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

//var TypeSmartModule = reflect.TypeOf((*Module)(nil)).Elem()

func init() {
	gs.Banner("")
	//gs.Banner(" .oooooo..o                                          .   \nd8P'    `Y8                                        .o8   \nY88bo.      ooo. .oo.  .oo.    .oooo.   oooo d8b .o888oo \n `\"Y8888o.  `888P\"Y88bP\"Y88b  `P  )88b  `888\"\"8P   888   \n     `\"Y88b  888   888   888   .oP\"888   888       888   \noo     .d8P  888   888   888  d8(  888   888       888 . \n8\"\"88888P'  o888o o888o o888o `Y888\"\"8o d888b      \"888\" ")
}
