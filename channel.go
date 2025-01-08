package mr_smart

import (
	"context"
	"encoding/binary"
	"github.com/cloudwego/netpoll"
	"github.com/pkg/errors"
	"github.com/ywengineer/mr.smart/codec"
	"github.com/ywengineer/mr.smart/utility"
	"go.uber.org/zap"
)

type SocketChannel struct {
	ctx       context.Context
	fd        int
	conn      netpoll.Connection
	codec     codec.Codec
	byteOrder binary.ByteOrder
	worker    Worker
	handlers  []ChannelHandler
}

// Send all data and event callback run in worker related SocketChannel
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

// LaterRun run task in worker related SocketChannel
func (h *SocketChannel) LaterRun(task func()) {
	h.worker.Run(h.ctx, task)
}

func (h *SocketChannel) Close() error {
	return h.conn.Close()
}

func (h *SocketChannel) send(data []byte) error {
	if h.conn == nil {
		return errors.New("SocketChannel is not initialized correctly")
	}
	writer := h.conn.Writer()
	defer writer.Flush()
	if _, err := writer.WriteBinary(data); err != nil {
		utility.DefaultLogger().Error("write data error", zap.Error(err))
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

func (h *SocketChannel) onMessageRead(ctx context.Context) error {
	if len(h.handlers) > 0 {
		for _, handler := range h.handlers {
			if err := handler.OnMessage(ctx, h); err != nil {
				return err
			}
		}
	}
	return nil
}
