package mr_smart

import (
	"context"
	"encoding/binary"
	"github.com/ywengineer/mr.smart/codec"
	"github.com/ywengineer/mr.smart/message"
)

type ChannelHandler interface {
	OnOpen(channel *SocketChannel)
	OnClose(channel *SocketChannel)
	OnMessage(ctx context.Context, channel *SocketChannel, msg *message.ProtocolMessage) error
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
