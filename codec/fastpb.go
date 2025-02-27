package codec

import (
	"errors"
	"github.com/cloudwego/fastpb"
	"github.com/ywengineer/smart/pkg"
)

var fpbc = &fastpbCodec{}

func Fastpb() Codec {
	return fpbc
}

func NewFastpbCodec() Codec {
	return &fastpbCodec{}
}

// fastpbCodec uses json marshaler and unmarshaler.
type fastpbCodec struct{}

// Encode encodes an object into slice of bytes.
func (c *fastpbCodec) Encode(i interface{}) ([]byte, error) {
	if v, ok := i.(fastpb.Writer); ok {
		buf := make([]byte, v.Size())
		v.FastWrite(buf)
		return buf, nil
	}
	return nil, errors.New("fastpb codec encode not support")
}

// Decode decodes an object from slice of bytes.
func (c *fastpbCodec) Decode(reader pkg.Reader, i interface{}) error {
	if v, ok := i.(fastpb.Reader); ok {
		bytes, err := readAll(reader)
		if err != nil {
			return err
		}
		_, err = fastpb.ReadMessage(bytes, int8(fastpb.SkipTypeCheck), v)
		return err
	}
	return errors.New("fastpb codec decode not support")
}
