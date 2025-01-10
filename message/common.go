package message

//var TypeReducible = reflect.TypeOf((*Reducible)(nil)).Elem()

type ProtocolId int16

const ProtocolMetaBytes = 8

const (
	Smart ProtocolId = iota
)

//type Reducible interface {
//	Reset()
//}
