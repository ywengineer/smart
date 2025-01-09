package smart

import (
	"context"
	"github.com/ywengineer/smart/message"
	"strconv"
)

type GateMessageHandler struct {
}

func (gsh *GateMessageHandler) OnOpen(channel *SocketChannel) {

}

func (gsh *GateMessageHandler) OnClose(channel *SocketChannel) {

}

func (gsh *GateMessageHandler) OnMessage(ctx context.Context, h *SocketChannel, m *message.ProtocolMessage) error {
	if m.Header == nil {
		m.Header = map[string]string{}
	}
	m.Header[FROM] = strconv.Itoa(h.GetFd())
	// load balance redirect to
	//bytes, _ := proto.Marshal(m)
	//

	return nil
}
