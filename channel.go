package mr_smart

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/netpoll"
	"github.com/ywengineer/mr.smart/codec"
	"go.uber.org/zap"
)

type SocketChannel struct {
	ctx       context.Context
	fd        int
	conn      netpoll.Connection
	codec     codec.Codec
	byteOrder binary.ByteOrder
	worker    gopool.Pool
	handlers  []ChannelHandler
}

// all data and event callback run in worker related SocketChannel
func (h *SocketChannel) Send(msg interface{}) error {
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

// run task in worker related SocketChannel
func (h *SocketChannel) LaterRun(task func()) {
	h.worker.CtxGo(h.ctx, task)
}

func (h *SocketChannel) Close() error {
	return h.conn.Close()
}

func (h *SocketChannel) send(data []byte) error {
	if _, err := h.conn.Writer().WriteBinary(data); err != nil {
		serverLogger.Error("write data error", zap.Error(err))
		return err
	}
	return nil
}

func (h *SocketChannel) onOpen() {
	h.LaterRun(func() {
		if len(h.handlers) > 0 {
			for _, handler := range h.handlers {
				handler.OnOpen(h)
			}
		}
	})
}

func (h *SocketChannel) onClose() {
	h.LaterRun(func() {
		if len(h.handlers) > 0 {
			for _, handler := range h.handlers {
				handler.OnClose(h)
			}
		}
	})
}

func (h *SocketChannel) onMessageRead() error {
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
			req := &Request{
				channel:     h,
				messageCode: msgCode,
				body:        bodyBytes,
			}
			h.LaterRun(func() {
				dispatchRequest(req)
			})
		}
	}
	return nil
}
