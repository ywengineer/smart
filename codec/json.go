package codec

import (
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
)

// JSONCodec uses json marshaler and unmarshaler.
type JSONCodec struct{}

// Encode encodes an object into slice of bytes.
func (c *JSONCodec) Encode(i interface{}) ([]byte, error) {
	return sonic.Marshal(i)
}

// Decode decodes an object from slice of bytes.
func (c *JSONCodec) Decode(data []byte, i interface{}) error {
	d := decoder.NewDecoder(string(data))
	d.UseNumber()
	return d.Decode(i)
}
