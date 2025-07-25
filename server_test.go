package smart

import (
	"context"
	"encoding/binary"
	loader2 "gitee.com/ywengineer/smart-kit/pkg/loaders"
	"gitee.com/ywengineer/smart-kit/pkg/nacos"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"gitee.com/ywengineer/smart/codec"
	"gitee.com/ywengineer/smart/message"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestServerWithNacos(t *testing.T) {
	//network, addr := "tcp", "127.0.0.1:12345"
	nc, err := nacos.NewNacosConfigClient(
		"192.168.44.128", 8848, "/nacos", 5000,
		"a7aabc24-17a7-4ac5-978f-6f933ce19dd4", "nacos", "nacos",
		"debug",
	)
	assert.Nil(t, err)
	//
	srv, err := NewSmartServer(loader2.NewDefaultNacosLoader(nc, "smart.gate.yaml", loader2.NewYamlDecoder()),
		WithByteOrder(func() binary.ByteOrder {
			return binary.LittleEndian
		}),
		WithCodec(func() codec.Codec {
			return codec.NewSmartCodec(binary.LittleEndian)
		}),
		AppendMessageHandler(func() MessageHandler { return NewGameMessageHandler() }),
	)
	// register game logic module
	err = RegisterModule(&TestModule{})
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}
	assert.Nil(t, err)
	// 4. start smart server
	ctx, err := srv.Serve(context.Background())
	assert.Nil(t, err)
	t.Log("smart nacos server was started.")
	//
	_srv := srv.(*defaultServer)
	runClient(t, ctx, _srv.conf.Network, _srv.conf.Address, _srv.initializers)
	//
	_ = <-utilk.WatchQuitSignal()
	// 5. wait smart server shutdown
	t.Logf("smart nacos server was stopped. %v", srv.Shutdown(ctx))
}

func TestGNetServer(t *testing.T) {
	// 1. load config: config.LoadConfig()
	// 2. register handler: RegisterModule
	// 3. create a smart server
	network, addr := "tcp", "127.0.0.1:12345"

	srv, err := NewGNetServer(
		loader2.NewValueLoader(&loader2.Conf{Network: network, Address: addr, WorkerLoadBalance: "rr"}),
		WithByteOrder(func() binary.ByteOrder {
			return binary.LittleEndian
		}),
		WithCodec(func() codec.Codec {
			return codec.NewSmartCodec(binary.LittleEndian)
		}),
		AppendMessageHandler(func() MessageHandler { return NewGameMessageHandler() }),
	)
	// register game logic module
	err = RegisterModule(&TestModule{})
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}
	assert.Nil(t, err)
	//
	srv.SetOnTick(func(ctx context.Context) time.Duration {
		return time.Second
	})
	// 4. start smart server
	ctx, err := srv.Serve(context.Background())
	assert.Nil(t, err)
	t.Log("smart server was started.")
	_srv := srv.(*gnetServer)
	//
	runClient(t, ctx, network, addr, _srv.initializers)
	//
	_ = <-utilk.WatchQuitSignal()
	t.Log("smart server will be stop.")
	err = srv.Shutdown(ctx)
	// 5. wait smart server shutdown
	t.Logf("smart server was stopped. %v", err)
}

func runClient(t *testing.T, ctx context.Context, network, addr string, initializers []ChannelInitializer) {
	go func() {
		time.Sleep(5 * time.Second)
		var err error
		//
		channel := NewAutoCloseSmartClient(ctx, network, addr, initializers)
		//
		if channel == nil {
			t.Error("dial failed")
			return
		}
		//
		if err = channel.Send(&message.ProtocolMessage{
			Seq:     1,
			Route:   1001,
			Header:  map[string]string{},
			Codec:   message.Codec_JSON,
			Payload: []byte(`{"ping":1001, "extra": "from client"}`),
		}); err != nil {
			t.Errorf("send 1001 failed. %v", err)
		}
		//
		if err = channel.Send(&message.ProtocolMessage{
			Seq:     2,
			Route:   1002,
			Header:  map[string]string{},
			Codec:   message.Codec_JSON,
			Payload: []byte(`{"ping":1002, "extra": "from client"}`),
		}); err != nil {
			t.Errorf("send 1002 failed. %v", err)
		}
		//
		if err = channel.Send(&message.ProtocolMessage{
			Seq:     3,
			Route:   1003,
			Header:  map[string]string{},
			Codec:   message.Codec_JSON,
			Payload: []byte(`{"ping":1003, "extra": "from client"}`),
		}); err != nil {
			t.Errorf("send 1003 failed. %v", err)
		}
		//
		if err = channel.Send(&message.ProtocolMessage{
			Seq:     4,
			Route:   1004,
			Header:  map[string]string{},
			Codec:   message.Codec_JSON,
			Payload: []byte(`{"ping":1004, "extra": "from client"}`),
		}); err != nil {
			t.Errorf("send 1004 failed. %v", err)
		}
	}()
}
