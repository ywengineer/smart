package smart

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/cloudwego/netpoll"
	"github.com/ywengineer/smart/codec"
	"github.com/ywengineer/smart/message"
	"github.com/ywengineer/smart/utility"
	"go.uber.org/zap"
	"sync"
)

type SocketChannel struct {
	ctx          context.Context
	fd           int
	conn         netpoll.Connection
	codec        codec.Codec
	byteOrder    binary.ByteOrder
	worker       Worker
	handlers     []ChannelHandler
	interceptors []MessageInterceptor
	msgHandlers  []MessageHandler
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

// GetFd returns id of channel
func (h *SocketChannel) GetFd() int {
	return h.fd
}

func (h *SocketChannel) onOpen() {
	if len(h.handlers) > 0 {
		h.LaterRun(func() {
			for _, handler := range h.handlers {
				handler.OnOpen(h)
			}
		})
	}
}

func (h *SocketChannel) onClose() {
	if len(h.handlers) > 0 {
		h.LaterRun(func() {
			for _, handler := range h.handlers {
				handler.OnClose(h)
			}
		})
	}
}

func (h *SocketChannel) onMessageRead(ctx context.Context) error {
	msg := protocolMessagePool.Get()
	err := h.codec.Decode(h.conn.Reader(), msg)
	// parameter type not match *message.ProtocolMessage
	// or pkg is too big
	if errors.Is(err, codec.ErrParamMessage) || errors.Is(err, codec.ErrTooBig) {
		_ = h.Close()
		return err
	} else if errors.Is(err, codec.ErrPkgNotFull) { // pkg not full, skip
		return nil
	} else { // decode success
		h.LaterRun(func(msg *message.ProtocolMessage) func() {
			return func() {
				defer protocolMessagePool.Put(msg)
				//
				ctx = context.WithValue(ctx, CtxKeyFromClient, h.GetFd())
				ctx = context.WithValue(ctx, CtxKeySeq, msg.GetSeq())
				ctx = context.WithValue(ctx, CtxKeyHeader, msg.GetHeader())
				//
				if len(h.interceptors) > 0 {
					for _, handler := range h.interceptors {
						if err := handler.BeforeInvoke(ctx, h, msg); err != nil {
							return
						}
					}
				}
				//
				if len(h.msgHandlers) > 0 {
					for _, handler := range h.msgHandlers {
						if err := handler.OnMessage(ctx, h, msg); err != nil {
							return
						}
					}
				}
				//
				if len(h.interceptors) > 0 {
					for deep := len(h.interceptors) - 1; deep >= 0; deep-- {
						if err := h.interceptors[deep].AfterInvoke(ctx, h, msg); err != nil {
							return
						}
					}
				}
			}
		}(msg.(*message.ProtocolMessage)))
	} //
	return nil
}

var protocolMessagePool = &sync.Pool{
	New: func() interface{} {
		return &message.ProtocolMessage{}
	},
}
