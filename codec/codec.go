package codec

import (
	"fmt"
	"github.com/ywengineer/smart/pkg"
	"reflect"
)

var b = &byteCodec{}

func Byte() Codec {
	return b
}

func NewByteCodec() Codec {
	return &byteCodec{}
}

func readAll(reader pkg.Reader) ([]byte, error) {
	return reader.ReadBinary(reader.Len())
}

// Codec defines the interface that decode/encode payload.
type Codec interface {
	Encode(i interface{}) ([]byte, error)
	Decode(reader pkg.Reader, i interface{}) error
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
func (c *byteCodec) Decode(reader pkg.Reader, i interface{}) error {
	if byes, err := readAll(reader); err != nil {
		return err
	} else {
		reflect.Indirect(reflect.ValueOf(i)).SetBytes(byes)
		return nil
	}
}
