package codec

import (
	"fmt"
	"github.com/cloudwego/netpoll"
	"google.golang.org/protobuf/proto"
)

var pc = &protoCodec{}

func Protobuf() Codec {
	return pc
}

func NewProtobufCodec() Codec {
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
func (c *protoCodec) Decode(reader netpoll.Reader, i interface{}) error {
	if m, ok := i.(proto.Message); ok {
		bytes, e := readAll(reader)
		if e != nil {
			return e
		}
		e = proto.Unmarshal(bytes, m)
		return e
	}
	return fmt.Errorf("%T is not a proto.Message", i)
}
