package smart

import (
	"context"
	"encoding/binary"
	"github.com/cloudwego/netpoll"
	"github.com/pkg/errors"
	"github.com/ywengineer/smart/codec"
	"github.com/ywengineer/smart/message"
	"github.com/ywengineer/smart/utility"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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

func (h *SocketChannel) SendSmart(msg *message.ProtocolMessage) error {
	return h.Send(msg)
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
		return h.send(data.Bytes())
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

	reader := h.conn.Reader()
	// 消息结构(len(4) + protocol(2) + compress(1) + flags(1) + payload(len))
	if reader.Len() < message.ProtocolMetaBytes {
		// data is not enough, wait next
		return nil
	}
	if data, err := reader.Peek(message.ProtocolMetaBytes); err != nil {
		// read data error
		utility.DefaultLogger().Error("read msg meta failed.", zap.Error(err))
		return err
	} else {
		pkgSize := int(h.byteOrder.Uint32(data[:4]))
		// TODO
		//_ = int(h.byteOrder.Uint32(data[4:6])) // protocol
		//_ = int(h.byteOrder.Uint32(data[6:7])) // compress
		//_ = int(h.byteOrder.Uint32(data[7:8])) // flags
		// pkg size reach max size, close it
		if pkgSize >= 65535 {
			utility.DefaultLogger().Error("msg is too big. connection will be close", zap.Int("size", pkgSize))
			return h.Close()
		}
		//
		if reader.Len() < pkgSize+message.ProtocolMetaBytes {
			// game msg body is not enough. wait
			return nil
		} else {
			err = reader.Skip(message.ProtocolMetaBytes)
			if err != nil {
				utility.DefaultLogger().Error("failed to skip meta size.", zap.Error(err))
				return err
			}
			pkg, err := reader.ReadBinary(pkgSize)
			if err != nil {
				utility.DefaultLogger().Error("failed to read protocol bytes.", zap.Error(err))
				return err
			}
			req := getRequest()
			if err = proto.Unmarshal(pkg, req); err != nil {
				utility.DefaultLogger().Error("failed to decode bytes to ProtocolMessage.", zap.Error(err))
				releaseRequest(req)
				return err
			}
			h.LaterRun(func(req *message.ProtocolMessage) func() {
				return func() {
					defer releaseRequest(req)
					//
					if len(h.handlers) > 0 {
						for _, handler := range h.handlers {
							if err := handler.OnMessage(ctx, h, req); err != nil {
								break
							}
						}
					}
				}
			}(req))
		}
	} //
	return nil
}
