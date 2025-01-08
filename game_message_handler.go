package mr_smart

import (
	"context"
	"github.com/ywengineer/mr.smart/utility"
	"go.uber.org/zap"
)

func NewGameMessageHandler() ChannelHandler {
	return &gameMessageHandler{byteLen: 4, byteCode: 4}
}

type gameMessageHandler struct {
	byteLen  int
	byteCode int
}

func (gsh *gameMessageHandler) OnOpen(channel *SocketChannel) {

}

func (gsh *gameMessageHandler) OnClose(channel *SocketChannel) {

}

func (gsh *gameMessageHandler) OnMessage(ctx context.Context, h *SocketChannel) error {
	reader := h.conn.Reader()
	// 消息结构(len(4) + code(4) + body(len - 4))
	if reader.Len() < gsh.byteLen {
		// data is not enough, wait next
		return nil
	}
	if data, err := reader.Peek(gsh.byteLen); err != nil {
		// read data error
		utility.DefaultLogger().Error("read game msg length failed.", zap.Error(err))
		return err
	} else {
		pkgSize := int(h.byteOrder.Uint32(data))
		if reader.Len() < pkgSize+gsh.byteLen {
			// game msg body is not enough. wait
			return nil
		} else {
			_ = reader.Skip(gsh.byteLen)
			pkg, _ := reader.Slice(pkgSize)
			codeBytes, _ := pkg.ReadBinary(gsh.byteCode)
			req := getRequest()
			req.messageCode = int(h.byteOrder.Uint32(codeBytes))
			req.body, _ = pkg.ReadBinary(pkgSize - gsh.byteCode)
			_ = pkg.Release()
			h.LaterRun(func() {
				defer releaseRequest(req)
				hManager.invokeHandler(ctx, h, req)
			})
		}
	} //
	return nil
}
