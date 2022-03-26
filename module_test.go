package mr_smart

import (
	"github.com/ywengineer/mr.smart/codec"
	"go.uber.org/zap"
	"testing"
)

type Req struct {
	Ping  int    `json:"ping"`
	Extra string `json:"extra"`
}

func (r *Req) Reset() {
	r.Ping = -1
	r.Extra = ""
}

type Res struct {
	Pong int `json:"pong"`
}

type TestModule struct {
}

func (m *TestModule) Name() string {
	return "TestModule"
}

func (m *TestModule) Handler1001(channel *SocketChannel, req *Req) {
	srvLogger.Info("Handler1001 invoked", zap.Any("req", *req))
}

func (m *TestModule) Handler1002(channel *SocketChannel, req *Req) {
	srvLogger.Info("Handler1002 invoked", zap.Any("req", *req))
}

func (m *TestModule) Handler1003(channel *SocketChannel, req *Req) {
	srvLogger.Info("Handler1003 invoked", zap.Any("req", *req))
}

func (m *TestModule) Handler1004(channel *SocketChannel, req *Req) *Res {
	srvLogger.Info("Handler1004 invoked", zap.Any("req", *req))
	return &Res{
		Pong: req.Ping,
	}
}

func TestRegisterModule(t *testing.T) {
	err := RegisterModule(&TestModule{})
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}
	jc := &codec.JSONCodec{}
	channel := &SocketChannel{codec: jc}
	channel.doRequest(&request{
		messageCode: 1001,
		body:        []byte(`{"ping": 1001, "extra": "1001"}`),
	})
	channel.doRequest(&request{
		messageCode: 1002,
		body:        []byte(`{"ping": 1002}`),
	})
	channel.doRequest(&request{
		messageCode: 1003,
		body:        []byte(`{"ping": 1003, "extra": "1003"}`),
	})
	channel.doRequest(&request{
		messageCode: 1004,
		body:        []byte(`{"ping": 1004}`),
	})
}
