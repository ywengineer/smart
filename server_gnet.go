package smart

import (
	"context"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/pkg/errors"
	"github.com/ywengineer/smart/pkg"
	"github.com/ywengineer/smart/utility"
	"go.uber.org/zap"
	"sync/atomic"
)

type gnetServer struct {
	baseServer
	gnet.BuiltinEventEngine
	eng          gnet.Engine
	disconnected int32
}

func (s *gnetServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	utility.DefaultLogger().Info("running server on " + fmt.Sprintf("%s://%s", s.conf.Network, s.conf.Address))
	s.eng = eng
	return
}

func (s *gnetServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	channel := &SocketChannel{
		ctx:  s.ctx,
		fd:   c.Fd(),
		conn: pkg.NetGNetConn(c),
	}
	s.onChannelOpen(channel)
	return
}

func (s *gnetServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		utility.DefaultLogger().Info("error occurred on channel", zap.String("remote", c.RemoteAddr().String()), zap.Error(err))
	}
	atomic.AddInt32(&s.disconnected, 1)
	_ = s.onChannelClosed(c.Fd())
	return
}

func (s *gnetServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	fd := c.Fd()
	if err := s.onChannelRead(s.ctx, fd); err != nil {
		if errors.Is(err, ErrNotRegisteredChannel) {
			utility.DefaultLogger().Error("not registered channel.", zap.Int("fd", fd))
			action = gnet.Close
		}
	}
	return
}

func (s *gnetServer) onSpin(ctx context.Context) error {
	return gnet.Run(s, s.conf.Network+"://"+s.conf.Address,
		gnet.WithMulticore(true),
		gnet.WithReusePort(true),
		gnet.WithReuseAddr(true),
		gnet.WithTCPNoDelay(gnet.TCPNoDelay),
		gnet.WithLogPath("./gnet.log"),
		gnet.WithLogLevel(utility.GetLogLevel()),
		gnet.WithSocketRecvBuffer(1*1024*1024),
	)
}

func (s *gnetServer) onShutdown() error {
	return s.eng.Stop(s.ctx)
}
