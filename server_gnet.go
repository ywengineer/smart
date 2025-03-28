package smart

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/pkg/errors"
	"github.com/ywengineer/smart-kit/pkg/logk"
	"github.com/ywengineer/smart/pkg"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

type gnetServer struct {
	*baseServer
	gnet.BuiltinEventEngine
	eng          gnet.Engine
	disconnected int32
}

func (s *gnetServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logk.Info("running server on " + fmt.Sprintf("%s://%s", s.conf.Network, s.conf.Address))
	s.eng = eng
	return
}

func (s *gnetServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	s.onChannelOpen(pkg.NetGNetConn(c))
	return
}

func (s *gnetServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logk.Info("error occurred on channel", zap.String("remote", c.RemoteAddr().String()), zap.Error(err))
	}
	atomic.AddInt32(&s.disconnected, 1)
	_ = s.onChannelClosed(c.Fd())
	return
}

func (s *gnetServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	fd := c.Fd()
	if err := s.onChannelRead(fd); err != nil {
		if errors.Is(err, ErrNotRegisteredChannel) {
			logk.Error("not registered channel.", zap.Int("fd", fd))
			action = gnet.Close
		}
	}
	return
}

func (s *gnetServer) onSpin() error {
	return gnet.Run(s, s.conf.Network+"://"+s.conf.Address,
		gnet.WithMulticore(true),
		gnet.WithTCPNoDelay(gnet.TCPNoDelay),
		gnet.WithLogger(logk.DefaultLogger()),
		gnet.WithSocketRecvBuffer(1*1024*1024),
		gnet.WithTCPKeepAlive(time.Minute), // keep alive
	)
}

func (s *gnetServer) onShutdown() error {
	return s.eng.Stop(s.ctx)
}
