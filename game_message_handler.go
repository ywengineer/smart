package smart

import (
	"context"
	"github.com/ywengineer/smart/message"
)

func NewGameMessageHandler() MessageHandler {
	return &gameMessageHandler{}
}

type gameMessageHandler struct {
}

func (gsh *gameMessageHandler) OnMessage(ctx context.Context, h Channel, m *message.ProtocolMessage) error {
	hManager.invokeHandler(ctx, h, m)
	return nil
}
