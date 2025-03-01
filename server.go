package smart

import (
	"context"
	"github.com/cloudwego/netpoll"
	"github.com/pkg/errors"
	sl "github.com/ywengineer/smart/loader"
	"github.com/ywengineer/smart/pkg"
	"github.com/ywengineer/smart/utility"
	"go.uber.org/zap"
	"net"
	"runtime"
	"time"
)

var ErrNotRegisteredChannel = errors.New("registered channel, close it.")

type status int32

const (
	prepared status = iota
	running
	stopped
)

type serverHolder interface {
	onSpin(ctx context.Context) error
	onShutdown() error
}

// Server smart server interface
type Server interface {
	Serve(ctx context.Context) (context.Context, error)
	Shutdown() error
	ConnCount() int32
	GetChannel(id int) (Channel, bool)
	SetOnConfigChange(callback func(conf sl.Conf))
	SetOnTick(tick func(ctx context.Context) time.Duration)
}

type defaultServer struct {
	*baseServer
	eventLoop netpoll.EventLoop
}

func _newServer(loader sl.SmartLoader, useGNet bool, initializer ...ChannelInitializer) (Server, error) {
	if len(initializer) == 0 {
		return nil, errors.New("initializer of channel can not be empty")
	}
	// load loader
	conf := &sl.Conf{}
	if err := loader.Load(conf); err != nil {
		return nil, errors.WithMessage(err, "load server loader error")
	} else {
		utility.DefaultLogger().Debug("new smart server with conf", zap.Any("conf", *conf))
	}
	worker, _ := NewWorkerManager(utility.MaxInt(conf.Workers, 1), parseLoadBalance(conf.WorkerLoadBalance))
	if useGNet {
		srv := &gnetServer{
			baseServer: &baseServer{
				status:        prepared,
				workerManager: worker,
				initializers:  initializer,
				conf:          conf,
				confWatcher:   loader.Watch,
			},
		}
		srv.baseServer.holder = srv
		return srv, nil
	} else {
		srv := &defaultServer{
			baseServer: &baseServer{
				status:        prepared,
				workerManager: worker,
				initializers:  initializer,
				conf:          conf,
				confWatcher:   loader.Watch,
			},
		}
		srv.baseServer.holder = srv
		return srv, nil
	}
}

func NewGNetServer(loader sl.SmartLoader, initializer ...ChannelInitializer) (Server, error) {
	return _newServer(loader, true, initializer...)
}

func NewSmartServer(loader sl.SmartLoader, initializer ...ChannelInitializer) (Server, error) {
	return _newServer(loader, false, initializer...)
}

func (s *defaultServer) onSpin(ctx context.Context) error {
	var listener net.Listener
	var err error
	if goos := runtime.GOOS; goos == "windows" {
		if listener, err = net.Listen(s.conf.Network, s.conf.Address); err != nil {
			utility.DefaultLogger().Panic("create server listener on windows error", zap.Error(err))
		}
	} else if listener, err = netpoll.CreateListener(s.conf.Network, s.conf.Address); err != nil {
		utility.DefaultLogger().Panic("create server listener error", zap.Error(err))
	}
	eventLoop, _ := netpoll.NewEventLoop(s.onConnRead, netpoll.WithOnPrepare(s.onConnPrepare), netpoll.WithOnConnect(s.onConnOpen))
	s.eventLoop = eventLoop
	//
	utility.DefaultLogger().Info("serve run at", zap.Any("address", s.conf.Network+"://"+s.conf.Address))
	return s.eventLoop.Serve(listener)
}

func (s *defaultServer) onShutdown() error {
	return s.eventLoop.Shutdown(s.ctx)
}

func (s *defaultServer) onConnPrepare(conn netpoll.Connection) context.Context {
	return s.ctx
}

func (s *defaultServer) onConnRead(_ context.Context, conn netpoll.Connection) error {
	fd := conn.(netpoll.Conn).Fd()
	if err := s.onChannelRead(fd); err != nil {
		if errors.Is(err, ErrNotRegisteredChannel) {
			utility.DefaultLogger().Error("not registered channel.", zap.Int("fd", fd))
			_ = conn.Close()
		}
		return err
	}
	return nil
}

func (s *defaultServer) onConnOpen(_ context.Context, conn netpoll.Connection) context.Context {
	_ = conn.AddCloseCallback(s.onConnClosed)
	s.onChannelOpen(pkg.NetNetpollConn(conn))
	return s.ctx
}

func (s *defaultServer) onConnClosed(conn netpoll.Connection) error {
	return s.onChannelClosed(conn.(netpoll.Conn).Fd())
}
