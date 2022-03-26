package mr_smart

import (
	"github.com/ywengineer/mr.smart/codec"
	"testing"
)

type Req struct {
	Ping  int    `json:"ping"`
	Extra string `json:"extra"`
}

type Res struct {
	Pong int `json:"pong"`
}

type Module struct {
}

type NotModule struct {
}

func (m *Module) Name() string {
	return "TestModule"
}

func (m *Module) Handler1001(channel *SocketChannel, req *Req) {
	srvLogger.Info("Handler1001 invoked")
}

func (m *Module) Handler1002(channel *SocketChannel, req *Req) {
	srvLogger.Info("Handler1002 invoked")
}

func (m *Module) Handler1003(channel *SocketChannel, req *Req) {
	srvLogger.Info("Handler1003 invoked")
}

func (m *Module) Handler1004(channel *SocketChannel, req *Req) *Res {
	srvLogger.Info("Handler1004 invoked")
	return &Res{
		Pong: req.Ping,
	}
}

func TestRegisterModule(t *testing.T) {
	err := RegisterModule(&Module{})
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}
	jc := &codec.JSONCodec{}
	channel := &SocketChannel{codec: jc}
	channel.doRequest(&request{
		messageCode: 1001,
		body:        []byte(`{"ping": 1001}`),
	})
	channel.doRequest(&request{
		messageCode: 1002,
		body:        []byte(`{"ping": 1002}`),
	})
	channel.doRequest(&request{
		messageCode: 1003,
		body:        []byte(`{"ping": 1003}`),
	})
	channel.doRequest(&request{
		messageCode: 1004,
		body:        []byte(`{"ping": 1004}`),
	})
}
