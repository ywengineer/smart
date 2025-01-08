package codec

import (
	"fmt"
	"google.golang.org/protobuf/proto"
)

func Protobuf() Codec {
	return &protoCodec{}
}

// protoCodec uses protobuf marshaler and unmarshaler.
type protoCodec struct{}

// Encode encodes an object into slice of bytes.
func (c *protoCodec) Encode(i interface{}) ([]byte, error) {
	if m, ok := i.(proto.Message); ok {
		return proto.Marshal(m)
	}
	return nil, fmt.Errorf("%T is not a proto.Message", i)
}

// Decode decodes an object from slice of bytes.
func (c *protoCodec) Decode(data []byte, i interface{}) error {
	if m, ok := i.(proto.Message); ok {
		return proto.Unmarshal(data, m)
	}
	return fmt.Errorf("%T is not a proto.Message", i)
}
