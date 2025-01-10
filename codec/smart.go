package codec

import (
	"encoding/binary"
	"errors"
	"github.com/cloudwego/netpoll"
	"github.com/ywengineer/smart/message"
	"github.com/ywengineer/smart/utility"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"sync"
)

var littleSmartCodec = NewSmartCodec(binary.LittleEndian)
var bigSmartCodec = NewSmartCodec(binary.BigEndian)
var ErrorTooBig = errors.New("msg too big")
var ErrorPkgNotFull = errors.New("pkg not full")
var ErrorParamMessage = errors.New("parameter is not a type [*message.ProtocolMessage]")

var smartMessagePool = &sync.Pool{
	New: func() interface{} {
		return &message.ProtocolMessage{}
	},
}

func NewPooledSmartMessage() *message.ProtocolMessage {
	return smartMessagePool.Get().(*message.ProtocolMessage)
}

func PutPooledSmartMessage(msg *message.ProtocolMessage) {
	if msg == nil {
		smartMessagePool.Put(msg)
	}
}

func LittleSmart() Codec {
	return littleSmartCodec
}

func BigSmart() Codec {
	return bigSmartCodec
}

func NewSmartCodec(odr binary.ByteOrder) Codec {
	return &smartCodec{odr: odr}
}

// jsonCodec uses json marshaler and unmarshaler.
type smartCodec struct {
	odr binary.ByteOrder
}

// Encode encodes an object into slice of bytes.
func (c *smartCodec) Encode(i interface{}) ([]byte, error) {
	if req, ok := i.(*message.ProtocolMessage); !ok {
		bytes, _ := proto.Marshal(req)
		buffer := netpoll.NewLinkBuffer(message.ProtocolMetaBytes + len(bytes))
		defer buffer.Release()
		_, _ = buffer.WriteBinary(utility.Int32ToBytes(c.odr, int32(len(bytes)))) // body len
		_, _ = buffer.WriteBinary(utility.Int16ToBytes(c.odr, int16(0)))          // protocol
		_ = buffer.WriteByte(0)                                                   // compress
		_ = buffer.WriteByte(0)                                                   // flags
		_, _ = buffer.WriteBinary(bytes)
		return buffer.Bytes(), nil
	}
	return nil, ErrorParamMessage
}

// Decode decodes an object from slice of bytes.
func (c *smartCodec) Decode(reader netpoll.Reader, i interface{}) error {
	dl := reader.Len()
	// 消息结构(len(4) + protocol(2) + compress(1) + flags(1) + payload(len))
	if dl < message.ProtocolMetaBytes {
		// data is not enough, wait next
		return ErrorPkgNotFull
	}
	data, _ := reader.Peek(message.ProtocolMetaBytes)
	pkgSize := int(c.odr.Uint32(data[:4])) // body len
	_ = int(c.odr.Uint16(data[4:6]))       // protocol
	_ = int(data[6])                       // compress
	_ = int(data[7])                       // flags
	// pkg size reach max size, close it
	if pkgSize >= 65535 {
		utility.DefaultLogger().Error("msg is too big.", zap.Int("size", pkgSize))
		return ErrorTooBig
	}
	//
	if dl < pkgSize+message.ProtocolMetaBytes {
		// game msg body is not enough. wait
		return ErrorPkgNotFull
	} else {
		_ = reader.Skip(message.ProtocolMetaBytes)
		pkg, _ := reader.ReadBinary(pkgSize)
		if req, ok := i.(*message.ProtocolMessage); !ok {
			return ErrorParamMessage
		} else {
			if err := proto.Unmarshal(pkg, req); err != nil {
				utility.DefaultLogger().Error("failed to decode bytes to ProtocolMessage.", zap.Error(err))
				return err
			}
			return nil
		}
	}
}
