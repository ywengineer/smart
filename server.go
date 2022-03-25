package mr_smart

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/netpoll"
	"github.com/ywengineer/mr.smart/codec"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)

const MsgSizeLength = 4
const MsgSizeCode = 4

type smartServer struct {
	lock         sync.Mutex
	running      int32
	ctx          context.Context
	shutdownHook context.CancelFunc
	eventLoop    netpoll.EventLoop
	channels     sync.Map // key=fd, value=connection
	channelCount int32    // accept counter
	initializers []ChannelInitializer
}

func NewSmartServer(initializer []ChannelInitializer) (*smartServer, error) {
	if len(initializer) == 0 {
		return nil, errors.New("holder initializer can not be empty")
	}
	server := &smartServer{
		running:      0,
		initializers: initializer,
	}
	return server, nil
}

func (s *smartServer) Serve(network, addr string) (context.Context, error) {
	defer s.lock.Unlock()
	s.lock.Lock()
	listener, err := netpoll.CreateListener(network, addr)
	if err != nil {
		return nil, err
	}
	rootCtx, cancel := context.WithCancel(context.Background())
	eventLoop, _ := netpoll.NewEventLoop(s.onConnRead, netpoll.WithOnPrepare(s.onConnPrepare), netpoll.WithOnConnect(s.onConnOpen))
	s.eventLoop = eventLoop
	s.ctx = rootCtx
	s.shutdownHook = cancel
	if atomic.CompareAndSwapInt32(&s.running, 0, 1) == false {
		return nil, errors.New("start smart server failed")
	}
	// start listen loop ...
	go func() {
		err = eventLoop.Serve(listener)
		// 启动失败
		if err != nil {
			_ = s.Shutdown()
		}
	}()
	return rootCtx, nil
}

func (s *smartServer) Shutdown() error {
	defer s.lock.Unlock()
	s.lock.Lock()
	running := atomic.LoadInt32(&s.running)
	if running == 0 {
		return errors.New("mr. smart server is not running")
	}
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
		serverLogger.Error("not registered channel.", zap.Int("fd", fd))
		_ = conn.Close()
		return fmt.Errorf("channel [%d] not registered", fd)
	} else { // registered
		return channel.(*socketChannel).onMessageRead()
	}
}

func (s *smartServer) onConnOpen(ctx context.Context, conn netpoll.Connection) context.Context {
	_ = conn.AddCloseCallback(s.onConnClosed)
	channel := &socketChannel{
		ctx:  ctx,
		fd:   conn.(netpoll.Conn).Fd(),
		conn: conn,
	}
	channel.worker = workers.Pick(channel.fd)
	for _, initializer := range s.initializers {
		initializer(channel)
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
		ch.(*socketChannel).onClose()
	}
	return nil
}

func (s *smartServer) ConnCount() int32 {
	return atomic.LoadInt32(&s.channelCount)
}

// GetChannel by fd(id)
func (s *smartServer) GetChannel(id int) *socketChannel {
	if ch, ok := s.channels.Load(id); ok {
		return ch.(*socketChannel)
	}
	return nil
}
