package codec

import (
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/cloudwego/netpoll"
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
func (c *jsonCodec) Encode(i interface{}) (*netpoll.LinkBuffer, error) {
	if buf, err := sonic.Marshal(i); err != nil {
		return nil, err
	} else {
		return newLinkBuffer(buf), nil
	}
}

// Decode decodes an object from slice of bytes.
func (c *jsonCodec) Decode(reader netpoll.Reader, i interface{}) error {
	if bytes, e := readAll(reader); e != nil {
		return e
	} else {
		d := decoder.NewDecoder(string(bytes))
		d.UseNumber()
		return d.Decode(i)
	}
}
