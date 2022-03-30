package mr_smart

import (
	"context"
	"encoding/binary"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/netpoll"
	"github.com/pkg/errors"
	"github.com/ywengineer/mr.smart/codec"
	"go.uber.org/zap"
)

const protocolLengthBytes = 4
const protocolCodeBytes = 4

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
	if h.conn == nil {
		return errors.New("SocketChannel is not initialized correctly")
	}
	writer := h.conn.Writer()
	defer writer.Flush()
	if _, err := writer.WriteBinary(data); err != nil {
		srvLogger.Error("write data error", zap.Error(err))
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
	if reader.Len() < protocolLengthBytes {
		srvLogger.Info("not enough data")
		return errors.New("not enough data")
	}
	if data, err := reader.Peek(protocolLengthBytes); err != nil {
		srvLogger.Error("read length failed.", zap.Error(err))
		return err
	} else {
		pkgSize := int(h.byteOrder.Uint32(data))
		if reader.Len() < pkgSize+protocolLengthBytes {
			srvLogger.Info("message body is not enough")
			return errors.New("message body is not enough")
		} else {
			_ = reader.Skip(protocolLengthBytes)
			pkg, _ := reader.Slice(pkgSize)
			codeBytes, _ := pkg.ReadBinary(protocolCodeBytes)
			req := getRequest()
			req.messageCode = int(h.byteOrder.Uint32(codeBytes))
			req.body, _ = pkg.ReadBinary(pkgSize - protocolCodeBytes)
			_ = pkg.Release()
			h.LaterRun(func() {
				h.doRequest(req)
			})
		}
	} //
	return nil
}

func (h *SocketChannel) doRequest(req *request) {
	defer releaseRequest(req)
	hd := hManager.findHandlerDefinition(req.messageCode)
	if hd == nil {
		srvLogger.Info("handler definition not found for message code", zap.Int("msgCode", req.messageCode))
		return
	}
	in := hd.getIn()
	// decode message
	if err := h.codec.Decode(req.body, in); err != nil {
		// decode failed. close channel
		srvLogger.Info("decode message error. suspicious channel, close it.", zap.Error(err))
		_ = h.Close()
		hd.releaseIn(in)
		return
	}
	response := hd.invoke(h, in)
	// oneway message
	if response == nil {
		return
	}
	// send response
	if err := h.Send(response); err != nil {
		srvLogger.Error("send response error", zap.Error(err))
	}
}
