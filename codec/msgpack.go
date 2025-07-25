package codec

import (
	"bytes"
	"gitee.com/ywengineer/smart/pkg"
	"github.com/vmihailenco/msgpack/v5"
)

var msgpackc = &msgpackCodec{}

func Msgpack() Codec {
	return msgpackc
}

func NewMsgpackCodec() Codec {
	return &msgpackCodec{}
}

// msgpackCodec uses messagepack marshaler and unmarshaler.
type msgpackCodec struct{}

// Encode encodes an object into slice of bytes.
func (c *msgpackCodec) Encode(i interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	// enc.UseJSONTag(true)
	err := enc.Encode(i)
	return buf.Bytes(), err
}

// Decode decodes an object from slice of bytes.
func (c *msgpackCodec) Decode(reader pkg.Reader, i interface{}) error {
	if buf, err := readAll(reader); err != nil {
		return err
	} else {
		dec := msgpack.NewDecoder(bytes.NewReader(buf))
		// dec.UseJSONTag(true)
		return dec.Decode(i)
	}
}
