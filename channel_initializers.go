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

func AddChannelHandler(f func() ChannelHandler) ChannelInitializer {
	return func(channel *SocketChannel) {
		channel.handlers = append(channel.handlers, f())
	}
}
