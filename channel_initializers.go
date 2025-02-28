package smart

import (
	"context"
	"encoding/binary"
	"github.com/ywengineer/smart/codec"
	"github.com/ywengineer/smart/message"
)

type ChannelHandler interface {
	OnOpen(channel Channel)
	OnClose(channel Channel)
}

type MessageHandler interface {
	// OnMessage skip execute next handler when return error
	OnMessage(ctx context.Context, channel Channel, msg *message.ProtocolMessage) error
}

// MessageInterceptor interceptor message handler
type MessageInterceptor interface {
	// BeforeInvoke skip execute next interceptor and message handler when return error
	BeforeInvoke(ctx context.Context, channel Channel, msg *message.ProtocolMessage) error
	// AfterInvoke skip execute next interceptor when return error
	AfterInvoke(ctx context.Context, channel Channel, msg *message.ProtocolMessage) error
}

type ChannelInitializer func(channel Channel)

func WithCodec(f func() codec.Codec) ChannelInitializer {
	return func(channel Channel) {
		channel.(*defaultChannel).codec = f()
	}
}

func WithByteOrder(f func() binary.ByteOrder) ChannelInitializer {
	return func(channel Channel) {
		channel.(*defaultChannel).byteOrder = f()
	}
}

func AppendHandler(f func() ChannelHandler) ChannelInitializer {
	return func(channel Channel) {
		channel.(*defaultChannel).handlers = append(channel.(*defaultChannel).handlers, f())
	}
}

func AppendMessageHandler(f func() MessageHandler) ChannelInitializer {
	return func(channel Channel) {
		channel.(*defaultChannel).msgHandlers = append(channel.(*defaultChannel).msgHandlers, f())
	}
}

func AppendMessageInterceptor(f func() MessageInterceptor) ChannelInitializer {
	return func(channel Channel) {
		channel.(*defaultChannel).interceptors = append(channel.(*defaultChannel).interceptors, f())
	}
}

func PrependHandler(f func() ChannelHandler) ChannelInitializer {
	return func(channel Channel) {
		channel.(*defaultChannel).handlers = append([]ChannelHandler{f()}, channel.(*defaultChannel).handlers...)
	}
}

func InsertHandlerAt(f func() ChannelHandler, pos int) ChannelInitializer {
	if pos <= 0 {
		return PrependHandler(f)
	}
	return func(channel Channel) {
		if pos <= 0 {
			PrependHandler(f)(channel)
		} else if pos >= len(channel.(*defaultChannel).handlers) {
			AppendHandler(f)(channel) // channel.handlers = append(channel.handlers, f())
		} else {
			hs := channel.(*defaultChannel).handlers[:pos]
			hs = append(hs, f())
			hs = append(hs, channel.(*defaultChannel).handlers[pos:]...)
			channel.(*defaultChannel).handlers = hs
		}
	}
}
