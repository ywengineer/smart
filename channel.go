package mr_smart

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/netpoll"
	"github.com/ywengineer/mr.smart/codec"
	"go.uber.org/zap"
)

type socketChannel struct {
	ctx       context.Context
	fd        int
	conn      netpoll.Connection
	codec     codec.Codec
	byteOrder binary.ByteOrder
	worker    gopool.Pool
	handlers  []ChannelHandler
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

func (h *socketChannel) onOpen() {
	if len(h.handlers) > 0 {
		for _, handler := range h.handlers {
			handler.OnOpen(h)
		}
	}
}

func (h *socketChannel) onClose() {
	if len(h.handlers) > 0 {
		for _, handler := range h.handlers {
			handler.OnClose(h)
		}
	}
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
			req := getRequestBodyModel(msgCode)
			if req == nil {
				serverLogger.Info("request type not found for message code", zap.Int("msgCode", msgCode))
				return fmt.Errorf("request type not found for message code: %d", msgCode)
			}
			if err = h.codec.Decode(bodyBytes, req); err != nil {
				serverLogger.Info("decode message error", zap.Error(err))
				return fmt.Errorf("decode message error: %s", err.Error())
			}
			dispatchRequest(&Request{
				request:     h,
				messageCode: msgCode,
				body:        req,
			})
		}
	}
	return nil
}
