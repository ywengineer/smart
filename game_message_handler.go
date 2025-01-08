package mr_smart

import (
	"context"
	"github.com/ywengineer/mr.smart/utility"
	"go.uber.org/zap"
)

const protocolLengthBytes = 4
const protocolCodeBytes = 4

type GameMessageHandler struct {
}

func (gsh *GameMessageHandler) OnOpen(channel *SocketChannel) {

}

func (gsh *GameMessageHandler) OnClose(channel *SocketChannel) {

}

func (gsh *GameMessageHandler) OnMessage(ctx context.Context, h *SocketChannel) error {
	reader := h.conn.Reader()
	// 消息结构(len(4) + code(4) + body(len - 4))
	if reader.Len() < protocolLengthBytes {
		// data is not enough, wait next
		return nil
	}
	if data, err := reader.Peek(protocolLengthBytes); err != nil {
		// read data error
		utility.DefaultLogger().Error("read game msg length failed.", zap.Error(err))
		return err
	} else {
		pkgSize := int(h.byteOrder.Uint32(data))
		if reader.Len() < pkgSize+protocolLengthBytes {
			// game msg body is not enough. wait
			return nil
		} else {
			_ = reader.Skip(protocolLengthBytes)
			pkg, _ := reader.Slice(pkgSize)
			codeBytes, _ := pkg.ReadBinary(protocolCodeBytes)
			req := getRequest()
			req.messageCode = int(h.byteOrder.Uint32(codeBytes))
			req.body, _ = pkg.ReadBinary(pkgSize - protocolCodeBytes)
			_ = pkg.Release()
			h.LaterRun(func() {
				defer releaseRequest(req)
				hManager.invokeHandler(ctx, h, req)
			})
		}
	} //
	return nil
}
