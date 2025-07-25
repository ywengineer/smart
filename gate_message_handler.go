package smart

import (
	"context"
	"gitee.com/ywengineer/smart/message"
	"strconv"
)

func NewGateMessageHandler() MessageHandler {
	return &gateMessageHandler{}
}

type gateMessageHandler struct {
}

func (gsh *gateMessageHandler) OnMessage(ctx context.Context, h Channel, m *message.ProtocolMessage) error {
	if m.Header == nil {
		m.Header = map[string]string{}
	}
	m.Header[HeaderFrom] = strconv.Itoa(h.GetFd())
	// load balance redirect to
	//bytes, _ := proto.Marshal(m)
	//

	return nil
}
