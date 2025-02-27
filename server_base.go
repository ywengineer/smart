package smart

import (
	"context"
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/ywengineer/smart/codec"
	sl "github.com/ywengineer/smart/loader"
	"github.com/ywengineer/smart/utility"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)

type baseServer struct {
	lock           sync.Mutex
	status         status
	channels       sync.Map // key=fd, value=connection
	channelCount   int32    // accept counter
	initializers   []ChannelInitializer
	workerManager  WorkerManager
	conf           *sl.Conf
	confWatcher    func(ctx context.Context, callback sl.WatchCallback) error
	onConfigChange func(conf sl.Conf)
	ctx            context.Context
	shutdownHook   context.CancelFunc
}

func (s *baseServer) onChannelRead(ctx context.Context, fd int) error {
	// channel not registered
	if channel, ok := s.channels.Load(fd); ok == false {
		return ErrNotRegisteredChannel
	} else { // registered
		return channel.(*SocketChannel).onMessageRead(ctx)
	}
}

func (s *baseServer) onChannelClosed(fd int) error {
	if ch, ok := s.channels.LoadAndDelete(fd); ok {
		atomic.AddInt32(&s.channelCount, -1)
		ch.(*SocketChannel).onClose()
	}
	return nil
}

func (s *baseServer) onChannelOpen(channel *SocketChannel) {
	channel.worker = s.workerManager.Pick(channel.fd)
	for _, initializer := range s.initializers {
		initializer(channel)
	}
	// check byte order
	if channel.byteOrder == nil {
		channel.byteOrder = binary.LittleEndian
		utility.DefaultLogger().Warn("byteOrder not set, default is LittleEndian")
	}
	// check codec, default codec
	if channel.codec == nil {
		channel.codec = codec.Byte()
		utility.DefaultLogger().Warn("codec not set, default is byte")
	}
	s.channels.Store(channel.fd, channel)
	atomic.AddInt32(&s.channelCount, 1)
	channel.onOpen()
}

func (s *baseServer) onSpin(ctx context.Context) error {
	return nil
}

func (s *baseServer) onShutdown() error {
	return nil
}

func (s *baseServer) ConnCount() int32 {
	return atomic.LoadInt32(&s.channelCount)
}

// GetChannel by fd(id)
func (s *baseServer) GetChannel(id int) *SocketChannel {
	if ch, ok := s.channels.Load(id); ok {
		return ch.(*SocketChannel)
	}
	return nil
}

func (s *baseServer) SetOnConfigChange(callback func(conf sl.Conf)) {
	s.onConfigChange = callback
}

func (s *baseServer) Shutdown() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.status != running {
		return errors.New("mr. smart server is not prepared or already stopped")
	}
	// can not start again.
	s.status = stopped
	s.shutdownHook()
	return s.onShutdown()
}

func (s *baseServer) Serve(ctx context.Context) (context.Context, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.status != prepared {
		return nil, errors.New("start smart server failed. maybe already started or can not start again")
	}
	s.ctx, s.shutdownHook = context.WithCancel(ctx)
	//
	s.status = running
	//
	s.ctx = context.WithValue(s.ctx, CtxKeyService, s.conf.ServiceName)
	// start listen loop ...
	go func() {
		//
		utility.DefaultLogger().Info("serve run at", zap.Any("address", s.conf.Network+s.conf.Address))
		//
		if err := s.onSpin(ctx); err != nil {
			utility.DefaultLogger().Panic("serve listener error", zap.Error(err))
			// start failed or serve quit
			_ = s.Shutdown()
		}
	}()
	// watch config
	if err := s.confWatcher(s.ctx, func(conf interface{}) error {
		utility.DefaultLogger().Debug("server config changed", zap.Any("old", *s.conf), zap.Any("new", *conf.(*sl.Conf)))
		s.conf = conf.(*sl.Conf)
		if s.onConfigChange != nil {
			s.onConfigChange(*s.conf)
		}
		return nil
	}); err != nil {
		utility.DefaultLogger().Error("server config watcher start error", zap.Error(err))
	} else {
		utility.DefaultLogger().Debug("server config watcher started")
	}
	return s.ctx, nil
}
