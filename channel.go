package smart

import (
	"context"
	"encoding/binary"
	"errors"
	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"gitee.com/ywengineer/smart/codec"
	"gitee.com/ywengineer/smart/message"
	"gitee.com/ywengineer/smart/pkg"
	"go.uber.org/zap"
	"sync"
)

type Channel interface {
	Context() context.Context
	LaterRun(task func())
	Close() error
	Send(msg interface{}) error
	GetFd() int
}

type defaultChannel struct {
	ctx          context.Context
	fd           int
	conn         pkg.Conn
	codec        codec.Codec
	byteOrder    binary.ByteOrder
	worker       Worker
	handlers     []ChannelHandler
	interceptors []MessageInterceptor
	msgHandlers  []MessageHandler
}

func (h *defaultChannel) Context() context.Context {
	return h.ctx
}

// Send all data and event callback run in worker related SocketChannel
func (h *defaultChannel) Send(msg interface{}) error {
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
func (h *defaultChannel) LaterRun(task func()) {
	h.worker.Run(h.ctx, task)
}

func (h *defaultChannel) Close() error {
	return h.conn.Close()
}

func (h *defaultChannel) send(data []byte) error {
	if h.conn == nil {
		return errors.New("SocketChannel is not initialized correctly")
	}
	writer := h.conn.Writer()
	defer writer.Flush()
	if _, err := writer.WriteBinary(data); err != nil {
		logk.Error("write data error", zap.Error(err))
		return err
	}
	return nil
}

// GetFd returns id of channel
func (h *defaultChannel) GetFd() int {
	return h.fd
}

func (h *defaultChannel) onOpen() {
	if len(h.handlers) > 0 {
		h.LaterRun(func() {
			for _, handler := range h.handlers {
				handler.OnOpen(h)
			}
		})
	}
}

func (h *defaultChannel) onClose() {
	if len(h.handlers) > 0 {
		h.LaterRun(func() {
			defer channelPool.Put(h)
			for _, handler := range h.handlers {
				handler.OnClose(h)
			}
		})
	}
}

func (h *defaultChannel) onMessageRead() error {
	for {
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
					ctx := context.WithValue(h.ctx, CtxKeySeq, msg.GetSeq())
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
	}
}

var protocolMessagePool = &sync.Pool{
	New: func() interface{} {
		return &message.ProtocolMessage{}
	},
}

var channelPool = &sync.Pool{
	New: func() interface{} {
		return &defaultChannel{}
	},
}
