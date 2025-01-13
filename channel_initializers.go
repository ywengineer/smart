package smart

import (
	"context"
	"encoding/binary"
	"github.com/ywengineer/smart/codec"
	"github.com/ywengineer/smart/message"
)

type ChannelHandler interface {
	OnOpen(channel *SocketChannel)
	OnClose(channel *SocketChannel)
}

type MessageHandler interface {
	// OnMessage skip execute next handler when return error
	OnMessage(ctx context.Context, channel *SocketChannel, msg *message.ProtocolMessage) error
}

// MessageInterceptor interceptor message handler
type MessageInterceptor interface {
	// BeforeInvoke skip execute next interceptor and message handler when return error
	BeforeInvoke(ctx context.Context, channel *SocketChannel, msg *message.ProtocolMessage) error
	// AfterInvoke skip execute next interceptor when return error
	AfterInvoke(ctx context.Context, channel *SocketChannel, msg *message.ProtocolMessage) error
}

type ChannelInitializer func(channel *SocketChannel)

func WithCodec(f func() codec.Codec) ChannelInitializer {
	return func(channel *SocketChannel) {
		channel.codec = f()
	}
}

func WithByteOrder(f func() binary.ByteOrder) ChannelInitializer {
	return func(channel *SocketChannel) {
		channel.byteOrder = f()
	}
}

func AppendHandler(f func() ChannelHandler) ChannelInitializer {
	return func(channel *SocketChannel) {
		channel.handlers = append(channel.handlers, f())
	}
}

func AppendMessageHandler(f func() MessageHandler) ChannelInitializer {
	return func(channel *SocketChannel) {
		channel.msgHandlers = append(channel.msgHandlers, f())
	}
}

func AppendMessageInterceptor(f func() MessageInterceptor) ChannelInitializer {
	return func(channel *SocketChannel) {
		channel.interceptors = append(channel.interceptors, f())
	}
}

func PrependHandler(f func() ChannelHandler) ChannelInitializer {
	return func(channel *SocketChannel) {
		channel.handlers = append([]ChannelHandler{f()}, channel.handlers...)
	}
}

func InsertHandlerAt(f func() ChannelHandler, pos int) ChannelInitializer {
	if pos <= 0 {
		return PrependHandler(f)
	}
	return func(channel *SocketChannel) {
		if pos <= 0 {
			PrependHandler(f)(channel)
		} else if pos >= len(channel.handlers) {
			AppendHandler(f)(channel) // channel.handlers = append(channel.handlers, f())
		} else {
			hs := channel.handlers[:pos]
			hs = append(hs, f())
			hs = append(hs, channel.handlers[pos:]...)
			channel.handlers = hs
		}
	}
}
