package mr_smart

import (
	"encoding/binary"
	"github.com/ywengineer/mr.smart/codec"
)

type ChannelHandler interface {
	OnOpen(channel *SocketChannel)
	OnClose(channel *SocketChannel)
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

func AddLastHandler(f func() ChannelHandler) ChannelInitializer {
	return func(channel *SocketChannel) {
		channel.handlers = append(channel.handlers, f())
	}
}

func AddFirstHandler(f func() ChannelHandler) ChannelInitializer {
	return func(channel *SocketChannel) {
		channel.handlers = append([]ChannelHandler{f()}, channel.handlers...)
	}
}

func AddHandlerAt(f func() ChannelHandler, pos int) ChannelInitializer {
	if pos <= 0 {
		return AddFirstHandler(f)
	}
	return func(channel *SocketChannel) {
		if pos <= 0 {
			AddFirstHandler(f)(channel)
		} else if pos >= len(channel.handlers) {
			AddLastHandler(f)(channel) // channel.handlers = append(channel.handlers, f())
		} else {
			hs := channel.handlers[:pos]
			hs = append(hs, f())
			hs = append(hs, channel.handlers[pos:]...)
			channel.handlers = hs
		}
	}
}
