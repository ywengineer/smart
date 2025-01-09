package message

import "reflect"

var TypeReducible = reflect.TypeOf((*Reducible)(nil)).Elem()

const ProtocolMetaBytes = 8

type Reducible interface {
	Reset()
}
