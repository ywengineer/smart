package mr_smart

import (
	"context"
	"github.com/ywengineer/mr.smart/message"
)

func NewGameMessageHandler() ChannelHandler {
	return &gameMessageHandler{}
}

type gameMessageHandler struct {
}

func (gsh *gameMessageHandler) OnOpen(channel *SocketChannel) {

}

func (gsh *gameMessageHandler) OnClose(channel *SocketChannel) {

}

func (gsh *gameMessageHandler) OnMessage(ctx context.Context, h *SocketChannel, m *message.ProtocolMessage) error {
	hManager.invokeHandler(ctx, h, m)
	return nil
}
