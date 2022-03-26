package message

import "reflect"

var TypeReducible = reflect.TypeOf((*Reducible)(nil)).Elem()

type Reducible interface {
	Reset()
}
