package codec

import (
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
)

var jsonc = &jsonCodec{}

func Json() Codec {
	return jsonc
}

func NewJsonCodec() Codec {
	return &jsonCodec{}
}

// jsonCodec uses json marshaler and unmarshaler.
type jsonCodec struct{}

// Encode encodes an object into slice of bytes.
func (c *jsonCodec) Encode(i interface{}) ([]byte, error) {
	return sonic.Marshal(i)
}

// Decode decodes an object from slice of bytes.
func (c *jsonCodec) Decode(data []byte, i interface{}) error {
	d := decoder.NewDecoder(string(data))
	d.UseNumber()
	return d.Decode(i)
}
