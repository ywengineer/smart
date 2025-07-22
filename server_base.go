package smart

import (
	"context"
	"encoding/binary"
	"gitee.com/ywengineer/smart/codec"
	"gitee.com/ywengineer/smart/pkg"
	"github.com/go-spring/spring-core/gs"
	"github.com/pkg/errors"
	"github.com/ywengineer/smart-kit/pkg/loaders"
	"github.com/ywengineer/smart-kit/pkg/logk"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

type baseServer struct {
	holder         serverHolder
	lock           sync.Mutex
	status         status
	channels       sync.Map // key=fd, value=connection
	channelCount   int32    // accept counter
	initializers   []ChannelInitializer
	workerManager  WorkerManager
	conf           *loaders.Conf
	confLoader     loaders.SmartLoader
	onConfigChange func(conf loaders.Conf)
	ctx            context.Context
	shutdownHook   context.CancelFunc
	onTick         func(ctx context.Context) time.Duration
}

func (s *baseServer) ticker() {
	if s.onTick == nil {
		return
	}
	var (
		delay time.Duration
		timer *time.Timer
	)
	defer func() {
		if timer != nil {
			timer.Stop()
		}
		if p := recover(); p != nil {
			if s.status == running {
				logk.Error("panic on server tick, restart it.", zap.Any("recover", p))
				go s.ticker()
			} else {
				logk.Error("panic on server tick", zap.Any("recover", p))
			}
		}
	}()
	for {
		delay = s.onTick(s.ctx)
		if timer == nil {
			timer = time.NewTimer(delay)
		} else {
			timer.Reset(delay)
		}
		select {
		case <-s.ctx.Done():
			break
		case <-timer.C:
		}
	}
}

func (s *baseServer) onChannelRead(fd int) error {
	// channel not registered
	if channel, ok := s.channels.Load(fd); ok == false {
		return ErrNotRegisteredChannel
	} else { // registered
		return channel.(*defaultChannel).onMessageRead()
	}
}

func (s *baseServer) onChannelClosed(fd int) error {
	if ch, ok := s.channels.LoadAndDelete(fd); ok {
		atomic.AddInt32(&s.channelCount, -1)
		ch.(*defaultChannel).onClose()
	}
	return nil
}

func (s *baseServer) onChannelOpen(conn pkg.Conn) {
	channel := channelPool.Get().(*defaultChannel)
	channel.ctx = context.WithValue(s.ctx, CtxKeyFromClient, conn.Fd())
	channel.conn, channel.fd = conn, conn.Fd()
	channel.worker = s.workerManager.Pick(channel.fd)
	for _, initializer := range s.initializers {
		initializer(channel)
	}
	// check byte order
	if channel.byteOrder == nil {
		channel.byteOrder = binary.LittleEndian
		logk.Warn("byteOrder not set, default is LittleEndian")
	}
	// check codec, default codec
	if channel.codec == nil {
		channel.codec = codec.Byte()
		logk.Warn("codec not set, default is byte")
	}
	s.channels.Store(channel.fd, channel)
	atomic.AddInt32(&s.channelCount, 1)
	channel.onOpen()
}

func (s *baseServer) ConnCount() int32 {
	return atomic.LoadInt32(&s.channelCount)
}

// GetChannel by fd(id)
func (s *baseServer) GetChannel(id int) (Channel, bool) {
	if ch, ok := s.channels.Load(id); ok {
		return ch.(Channel), ok
	}
	return nil, false
}

func (s *baseServer) SetOnConfigChange(callback func(conf loaders.Conf)) {
	s.onConfigChange = callback
}

func (s *baseServer) SetOnTick(tick func(ctx context.Context) time.Duration) {
	s.onTick = tick
}

func (s *baseServer) Shutdown(ctx context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.status != running {
		return errors.New("mr. smart server is not prepared or already stopped")
	}
	// can not start again.
	s.status = stopped
	s.shutdownHook()
	return s.holder.onShutdown()
}

func (s *baseServer) ListenAndServe(sig gs.ReadySignal) error {
	// start smart server
	_, err := s.Serve(context.Background())
	<-sig.TriggerAndWait()
	return err
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
		logk.Info("serve run at", zap.Any("address", s.conf.Network+"://"+s.conf.Address))
		//
		if err := s.holder.onSpin(); err != nil {
			logk.Fatal("serve listener error", zap.Error(err))
			// start failed or serve quit
			_ = s.Shutdown(s.ctx)
		}
	}()
	// tick
	go s.ticker()
	// watch config
	if err := s.confLoader.Watch(s.ctx, func(conf string) error {
		if err := s.confLoader.Unmarshal([]byte(conf), s.conf); err != nil {
			logk.Error("unmarshal configuration when watch", zap.Error(err))
			return err
		}
		logk.Debug("server config changed", zap.Any("conf", *s.conf))
		if s.onConfigChange != nil {
			s.onConfigChange(*s.conf)
		}
		return nil
	}); err != nil {
		logk.Error("server config watcher start error", zap.Error(err))
	} else {
		logk.Debug("server config watcher started")
	}
	return s.ctx, nil
}
