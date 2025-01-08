package codec

import (
	"fmt"
	"reflect"
)

func Byte() Codec {
	return &byteCodec{}
}

// Codec defines the interface that decode/encode payload.
type Codec interface {
	Encode(i interface{}) ([]byte, error)
	Decode(data []byte, i interface{}) error
}

// byteCodec uses raw slice pf bytes and don't encode/decode.
type byteCodec struct{}

// Encode returns raw slice of bytes.
func (c *byteCodec) Encode(i interface{}) ([]byte, error) {
	if data, ok := i.([]byte); ok {
		return data, nil
	}
	if data, ok := i.(*[]byte); ok {
		return *data, nil
	}

	return nil, fmt.Errorf("%T is not a []byte", i)
}

// Decode returns raw slice of bytes.
func (c *byteCodec) Decode(data []byte, i interface{}) error {
	reflect.Indirect(reflect.ValueOf(i)).SetBytes(data)
	return nil
}
