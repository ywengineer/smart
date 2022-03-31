package mr_smart

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"github.com/ywengineer/g-util/util"
	"github.com/ywengineer/mr.smart/codec"
	"github.com/ywengineer/mr.smart/server_config"
	"testing"
)

func TestServer(t *testing.T) {
	// 1. load config: config.LoadConfig()
	// 2. register handler: RegisterModule
	// 3. create a smart server
	srv, err := NewSmartServer(&server_config.ValueLoader{
		Conf: &server_config.Conf{Network: "tcp", Address: ":12345", Workers: 1, WorkerLoadBalance: "rr"},
	}, []ChannelInitializer{
		WithByteOrder(func() binary.ByteOrder {
			return binary.LittleEndian
		}),
		WithCodec(func() codec.Codec {
			return &codec.JSONCodec{}
		}),
	})
	assert.Nil(t, err)
	// 4. start smart server
	_, err = srv.Serve()
	assert.Nil(t, err)
	t.Log("smart server was started.")
	_ = <-util.WatchQuitSignal()
	// 5. wait smart server shutdown
	t.Logf("smart server was stopped. %v", srv.Shutdown())
}
