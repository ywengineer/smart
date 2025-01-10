package smart

import (
	"context"
	"encoding/binary"
	"github.com/cloudwego/netpoll"
	"github.com/stretchr/testify/assert"
	"github.com/ywengineer/smart/codec"
	"github.com/ywengineer/smart/message"
	"github.com/ywengineer/smart/server_config"
	"github.com/ywengineer/smart/utility"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	// 1. load config: config.LoadConfig()
	// 2. register handler: RegisterModule
	// 3. create a smart server
	network, addr := "tcp", "127.0.0.1:12345"
	srv, err := NewSmartServer(&server_config.ValueLoader{
		Conf: &server_config.Conf{Network: network, Address: addr, Workers: 1, WorkerLoadBalance: "rr"},
	},
		WithByteOrder(func() binary.ByteOrder {
			return binary.LittleEndian
		}),
		WithCodec(func() codec.Codec {
			return codec.NewSmartCodec(binary.LittleEndian)
		}),
		AppendHandler(func() ChannelHandler { return NewGameMessageHandler() }),
	)
	// register game logic module
	err = RegisterModule(&TestModule{})
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}
	assert.Nil(t, err)
	// 4. start smart server
	_, err = srv.Serve()
	assert.Nil(t, err)
	t.Log("smart server was started.")
	//
	dialer := netpoll.NewDialer()
	go func() {
		time.Sleep(5 * time.Second)
		conn, err := dialer.DialConnection(network, addr, time.Second)
		if err != nil {
			t.Errorf("dial failed. %v", err)
		}
		//
		channel := &SocketChannel{
			ctx:    context.Background(),
			fd:     conn.(netpoll.Conn).Fd(),
			conn:   conn,
			worker: srv.workerManager.Pick(conn.(netpoll.Conn).Fd()),
		}
		//conn.AddCloseCallback(channel.onClose)
		for _, initializer := range srv.initializers {
			initializer(channel)
		}
		channel.onOpen()
		defer channel.Close()
		//
		if err = channel.Send(&message.ProtocolMessage{
			Seq:     1,
			Route:   1001,
			Header:  map[string]string{},
			Codec:   message.Codec_JSON,
			Payload: []byte(`{"key":"1001"}`),
		}); err != nil {
			t.Errorf("send 1001 failed. %v", err)
		}
		//
		if err = channel.Send(&message.ProtocolMessage{
			Seq:     2,
			Route:   1001,
			Header:  map[string]string{},
			Codec:   message.Codec_JSON,
			Payload: []byte(`{"key":"1002"}`),
		}); err != nil {
			t.Errorf("send 1002 failed. %v", err)
		}
		//
		if err = channel.Send(&message.ProtocolMessage{
			Seq:     3,
			Route:   1001,
			Header:  map[string]string{},
			Codec:   message.Codec_JSON,
			Payload: []byte(`{"key":"1003"}`),
		}); err != nil {
			t.Errorf("send 1003 failed. %v", err)
		}
	}()
	_ = <-utility.WatchQuitSignal()
	// 5. wait smart server shutdown
	t.Logf("smart server was stopped. %v", srv.Shutdown())
}
