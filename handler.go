package smart

import (
	"context"
	"gitee.com/ywengineer/smart/codec"
	"gitee.com/ywengineer/smart/message"
	"github.com/ywengineer/smart-kit/pkg/logk"
	"github.com/ywengineer/smart-kit/pkg/utilk"
	"go.uber.org/zap"
	"reflect"
	"sync"
)

var hManager = &handlerManager{
	_handlerMap: make(map[int32]*handlerDefinition, 1000),
}

// handler structure
// code : handler for message code
// name : string
// in(context.Context, Channel, request): request must be a ptr
type handlerDefinition struct {
	messageCode int
	name        string
	method      reflect.Value
	inType      reflect.Type // must be ptr
	inPool      *sync.Pool
}

func (hd *handlerDefinition) invoke(ctx context.Context, channel Channel, in interface{}) interface{} {
	out := hd.method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(channel), reflect.ValueOf(in)})
	if len(out) == 0 {
		return nil
	}
	return out[0].Interface()
}

func (hd *handlerDefinition) releaseIn(in interface{}) {
	// release in to object pool
	hd.inPool.Put(in)
}

func (hd *handlerDefinition) newIn() interface{} {
	return hd.inPool.Get()
}

type handlerManager struct {
	// can not use directly
	_handlerMap map[int32]*handlerDefinition
}

func (hm *handlerManager) invokeHandler(ctx context.Context, c Channel, req *message.ProtocolMessage) {
	hd := hm.findHandlerDefinition(req.GetRoute())
	if hd == nil {
		logk.Error("handler definition not found for message code", zap.Int32("msgCode", req.GetRoute()))
		return
	}
	// find codec
	_codec := findMessageCodec(c, req.Codec)
	if _codec == nil {
		logk.Error("message codec not found", zap.String("codec", req.GetCodec().String()))
		_ = c.Close()
		return
	}
	in, buf := hd.newIn(), utilk.NewLinkBuffer(req.Payload)
	defer func() {
		hd.releaseIn(in)
		_ = buf.Release()
	}()
	// decode message
	if err := _codec.Decode(buf, in); err != nil {
		// decode failed. close channel
		logk.Error("decode message error. suspicious channel, close it.", zap.Error(err))
		_ = c.Close()
	} else if response := hd.invoke(ctx, c, in); response != nil {
		if err = c.Send(response); err != nil { // send response
			logk.Error("send response error", zap.Error(err))
		}
	} else { // oneway message
		// ignore
	}
}

func (hm *handlerManager) findHandlerDefinition(msgCode int32) *handlerDefinition {
	return hm._handlerMap[msgCode]
}

func (hm *handlerManager) addHandlerDefinition(def *handlerDefinition) {
	if _, ok := hm._handlerMap[int32(def.messageCode)]; ok {
		logk.Warnf("handler for message code [%d] already exists", def.messageCode)
	} else {
		logk.Debugf("register a new method handler for message code: %d", def.messageCode)
		def.inPool = &sync.Pool{
			New: func(hd *handlerDefinition) func() interface{} {
				return func() interface{} {
					in := hd.inType
					if in.Kind() == reflect.Ptr {
						in = in.Elem()
					}
					return reflect.New(in).Interface()
				}
			}(def),
		}
		hm._handlerMap[int32(def.messageCode)] = def
	}
}

func findMessageCodec(sc Channel, mc message.Codec) codec.Codec {
	switch mc {
	case message.Codec_JSON:
		return codec.Json()
	case message.Codec_PROTO:
		return codec.Protobuf()
	case message.Codec_MSGPACK:
		return codec.Msgpack()
	case message.Codec_THRIFT:
		logk.Warn("unsupported message codec: THRIFT")
		return nil
	case message.Codec_FAST_PB:
		return codec.Fastpb()
	case message.Codec_SERVER:
		return sc.(*defaultChannel).codec
	default:
		return nil
	}
}
