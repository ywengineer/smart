package smart

import (
	"context"
	"gitee.com/ywengineer/smart/pkg"
	"github.com/cloudwego/netpoll"
	"github.com/ywengineer/smart-kit/pkg/logk"
	"go.uber.org/zap"
	"strconv"
	"sync/atomic"
	"time"
)

var seqSmartClient uint64 = 1

func NewAutoCloseSmartClient(ctx context.Context, network, addr string, initializers []ChannelInitializer) Channel {
	return NewSmartClient(ctx, network, addr, initializers, true)
}

func NewSmartClient(ctx context.Context, network, addr string, initializers []ChannelInitializer, autoClose bool) Channel {
	//
	dialer := netpoll.NewDialer()
	//
	conn, err := dialer.DialConnection(network, addr, time.Second)
	//
	if err != nil {
		logk.Fatal("connect to smart server failed", zap.String("server", network+"://"+addr), zap.Error(err))
		return nil
	}
	//------------------------------------------------------------------------------------
	scId := strconv.FormatUint(atomic.AddUint64(&seqSmartClient, 1), 10)
	channel := &defaultChannel{
		ctx:  context.WithValue(ctx, CtxKeyFromClient, conn.(netpoll.Conn).Fd()),
		fd:   conn.(netpoll.Conn).Fd(),
		conn: pkg.NetNetpollConn(conn),
		worker: NewSingleWorker("smart-client-"+scId, func(ctx context.Context, i interface{}) {
			logk.Error("client worker panic occurred", zap.String("smart-client", scId), zap.Any("err", i))
		}),
	}
	for _, initializer := range initializers {
		initializer(channel)
	}
	//------------------------------------------------------------------------------------
	_ = conn.SetOnRequest(func(ctx context.Context, connection netpoll.Connection) error {
		return channel.onMessageRead()
	})
	_ = conn.AddCloseCallback(func(connection netpoll.Connection) error {
		channel.onClose()
		return nil
	})
	//------------------------------------------------------------------------------------
	channel.onOpen()
	// 自动关闭
	if autoClose {
		go func() {
			for {
				select {
				case <-ctx.Done():
					logk.Info(
						"client will be close, because of client running context is finished",
						zap.String("client", scId),
						zap.Error(channel.Close()))
					return
				default:
					if ctx.Err() != nil {
						logk.Error(
							"client will be close, because of client context error occurred",
							zap.String("client", scId),
							zap.Error(channel.Close()))
						return
					}
				}
			}
		}()
	}
	//
	return channel
}
