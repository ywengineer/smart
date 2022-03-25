package mr_smart

import (
	"encoding/binary"
	"github.com/ywengineer/mr.smart/codec"
)

type ChannelInitializer func(channel *socketChannel)

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
