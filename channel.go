package mr_smart

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/cloudwego/netpoll"
	"github.com/ywengineer/mr.smart/codec"
	"go.uber.org/zap"
)

type socketChannel struct {
	fd        int
	conn      netpoll.Connection
	codec     codec.Codec
	byteOrder binary.ByteOrder
}

func (h *socketChannel) Send(msg interface{}) error {
	// already encoded, send directly.
	if data, ok := msg.([]byte); ok {
		return h.send(data)
	} else if data, ok := msg.(*[]byte); ok {
		return h.send(*data)
	} else if data, err := h.codec.Encode(msg); err != nil { // encode message
		return err
	} else {
		return h.send(data)
	}
}

func (h *socketChannel) send(data []byte) error {
	_, err := h.conn.Writer().WriteBinary(data)
	return err
}

func (h *socketChannel) onMessageRead() error {
	reader := h.conn.Reader()
	// 消息结构(len(4) + code(4) + body(len - 4))
	if reader.Len() < MsgSizeLength {
		serverLogger.Info("not enough data")
		return errors.New("not enough data")
	}
	if data, err := reader.Peek(MsgSizeLength); err != nil {
		serverLogger.Error("read length failed.", zap.Error(err))
		return err
	} else {
		pkgSize := int(h.byteOrder.Uint32(data))
		if reader.Len() < pkgSize+MsgSizeLength {
			serverLogger.Info("message body is not enough")
			return errors.New("message body is not enough")
		} else {
			_ = reader.Skip(MsgSizeLength)
			pkg, _ := reader.Slice(pkgSize)
			defer pkg.Release()
			codeBytes, _ := pkg.Next(MsgSizeCode)
			msgCode := int(h.byteOrder.Uint32(codeBytes))
			bodyBytes, _ := pkg.ReadBinary(pkgSize - MsgSizeCode)
			if err = h.codec.Decode(bodyBytes, nil); err != nil {
				serverLogger.Info("decode message error", zap.Error(err))
				return fmt.Errorf("decode message error: %s", err.Error())
			}
		}
	}
	return nil
}
