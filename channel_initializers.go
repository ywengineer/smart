package mr_smart

import (
	"encoding/binary"
	"github.com/ywengineer/mr.smart/codec"
)

type ChannelHandler interface {
	OnOpen(channel *socketChannel)
	OnClose(channel *socketChannel)
}

type ChannelInitializer func(channel *socketChannel)
type OnOpen func(channel *socketChannel)

func SetCodec(f func() codec.Codec) ChannelInitializer {
	return func(channel *socketChannel) {
		channel.codec = f()
	}
}

func SetByteOrder(f func() binary.ByteOrder) ChannelInitializer {
	return func(channel *socketChannel) {
		channel.byteOrder = f()
	}
}

func AddChannelHandler(f func() ChannelHandler) ChannelInitializer {
	return func(channel *socketChannel) {
		channel.handlers = append(channel.handlers, f())
	}
}
