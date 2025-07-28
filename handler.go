package smart

import (
	"context"
	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"gitee.com/ywengineer/smart/codec"
	"gitee.com/ywengineer/smart/message"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"reflect"
	"sync"
)

var hManager = &handlerManager{
	_handlerMap: make(map[int32]*handlerDefinition, 1000),
}

type HandlerOutType int

const (
	HandlerOutTypeNil HandlerOutType = iota
	HandlerOutTypeByteSlice
	HandlerOutTypeProtoMessage
	HandlerOutTypeSmart
)

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
	outType     HandlerOutType
}

func (hd *handlerDefinition) invoke(ctx context.Context, channel Channel, in interface{}) (interface{}, interface{}) {
	out := hd.method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(channel), reflect.ValueOf(in)})
	if len(out) == 0 {
		return nil, nil
	} else if len(out) == 1 {
		return out[0].Interface(), nil
	} else {
		return out[0].Interface(), out[1].Interface()
	}
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
	} else if out0, out1 := hd.invoke(ctx, c, in); out0 != nil || out1 != nil {
		res := req
		if hd.outType == HandlerOutTypeProtoMessage {
			res.Route = int32(out0.(int))
			res.Payload, err = proto.Marshal(out1.(proto.Message))
			if err != nil {
				logk.Errorf("encode handler response error. route = %d, err = %v", hd.messageCode, err)
				return
			}
		} else if hd.outType == HandlerOutTypeByteSlice {
			res.Route = int32(out0.(int))
			res.Payload = out1.([]byte)
		} else if hd.outType == HandlerOutTypeSmart {
			res = out0.(*message.ProtocolMessage)
		} else {
			return
		}
		//
		if err = c.Send(res); err != nil { // send response
			logk.Errorf("send response error: %v", err)
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
