package mr_smart

import (
	"encoding/binary"
	"github.com/ywengineer/g-util/util"
	"github.com/ywengineer/mr.smart/codec"
	"testing"
)

func TestServer(t *testing.T) {
	srv, err := NewSmartServer([]ChannelInitializer{
		WithByteOrder(func() binary.ByteOrder {
			return binary.LittleEndian
		}),
		WithCodec(func() codec.Codec {
			return &codec.JSONCodec{}
		}),
	})
	if err != nil {
		t.Errorf("create smart server failed. %v", err)
		t.FailNow()
	}
	_, err = srv.Serve("tcp", ":12345")
	if err != nil {
		t.Errorf("start smart server failed. %v", err)
		t.FailNow()
	}
	t.Log("smart server was started.")
	_ = <-util.WatchQuitSignal()

	t.Logf("smart server was stopped. %v", srv.Shutdown())
}
