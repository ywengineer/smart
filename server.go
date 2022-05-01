package mr_smart

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/cloudwego/netpoll"
	"github.com/pkg/errors"
	"github.com/ywengineer/g-util/util"
	"github.com/ywengineer/mr.smart/codec"
	"github.com/ywengineer/mr.smart/log"
	"github.com/ywengineer/mr.smart/server_config"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)

type status int32

const (
	prepared status = iota
	running
	stopped
)

type smartServer struct {
	lock          sync.Mutex
	status        status
	ctx           context.Context
	shutdownHook  context.CancelFunc
	eventLoop     netpoll.EventLoop
	channels      sync.Map // key=fd, value=connection
	channelCount  int32    // accept counter
	initializers  []ChannelInitializer
	workerManager WorkerManager
	conf          *server_config.Conf
}

func NewSmartServer(loader server_config.Loader, initializer []ChannelInitializer) (*smartServer, error) {
	if len(initializer) == 0 {
		return nil, errors.New("initializer of channel can not be empty")
	}
	// load server_config
	conf, err := loader.Load()
	if err != nil || conf == nil {
		return nil, errors.WithMessage(err, "load server server_config error")
	}
	worker, _ := NewWorkerManager(util.MaxInt(conf.Workers, 1), parseLoadBalance(conf.WorkerLoadBalance))
	server := &smartServer{
		status:        prepared,
		workerManager: worker,
		initializers:  initializer,
		conf:          conf,
	}
	return server, nil
}

func (s *smartServer) Serve() (context.Context, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.status != prepared {
		return nil, errors.New("start smart server failed. maybe already started or can not start again")
	}
	listener, err := netpoll.CreateListener(s.conf.Network, s.conf.Address)
	if err != nil {
		return nil, err
	}
	s.ctx, s.shutdownHook = context.WithCancel(context.Background())
	eventLoop, _ := netpoll.NewEventLoop(s.onConnRead, netpoll.WithOnPrepare(s.onConnPrepare), netpoll.WithOnConnect(s.onConnOpen))
	s.eventLoop = eventLoop
	s.status = running
	// start listen loop ...
	go func() {
		err = eventLoop.Serve(listener)
		// start failed or serve quit
		if err != nil {
			_ = s.Shutdown()
		}
	}()
	return s.ctx, nil
}

func (s *smartServer) Shutdown() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.status != running {
		return errors.New("mr. smart server is not prepared or already stopped")
	}
	// can not start again.
	s.status = stopped
	s.shutdownHook()
	return s.eventLoop.Shutdown(s.ctx)
}

func (s *smartServer) onConnPrepare(conn netpoll.Connection) context.Context {
	return s.ctx
}

func (s *smartServer) onConnRead(ctx context.Context, conn netpoll.Connection) error {
	fd := conn.(netpoll.Conn).Fd()
	// channel not registered
	if channel, ok := s.channels.Load(fd); ok == false {
		log.GetLogger().Error("not registered channel.", zap.Int("fd", fd))
		_ = conn.Close()
		return fmt.Errorf("channel [%d] not registered", fd)
	} else { // registered
		return channel.(*SocketChannel).onMessageRead()
	}
}

func (s *smartServer) onConnOpen(ctx context.Context, conn netpoll.Connection) context.Context {
	_ = conn.AddCloseCallback(s.onConnClosed)
	channel := &SocketChannel{
		ctx:  ctx,
		fd:   conn.(netpoll.Conn).Fd(),
		conn: conn,
	}
	channel.worker = s.workerManager.Pick(channel.fd)
	for _, initializer := range s.initializers {
		initializer(channel)
	}
	// check byte order
	if channel.byteOrder == nil {
		channel.byteOrder = binary.LittleEndian
	}
	// check codec, default codec
	if channel.codec == nil {
		channel.codec = &codec.ByteCodec{}
	}
	s.channels.Store(channel.fd, channel)
	atomic.AddInt32(&s.channelCount, 1)
	channel.onOpen()
	return s.ctx
}

func (s *smartServer) onConnClosed(conn netpoll.Connection) error {
	fd := conn.(netpoll.Conn).Fd()
	atomic.AddInt32(&s.channelCount, -1)
	if ch, ok := s.channels.LoadAndDelete(fd); ok {
		ch.(*SocketChannel).onClose()
	}
	return nil
}

func (s *smartServer) ConnCount() int32 {
	return atomic.LoadInt32(&s.channelCount)
}

// GetChannel by fd(id)
func (s *smartServer) GetChannel(id int) *SocketChannel {
	if ch, ok := s.channels.Load(id); ok {
		return ch.(*SocketChannel)
	}
	return nil
}
