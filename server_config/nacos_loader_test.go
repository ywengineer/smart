package server_config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/ywengineer/smart/utility"
	"testing"
)

func TestNacosLoader(t *testing.T) {
	nc, err := NewNacosClient("192.168.44.128", 8848, "/nacos", 5000,
		"a7aabc24-17a7-4ac5-978f-6f933ce19dd4", "nacos", "nacos", "debug")
	assert.Nil(t, err)
	//
	loader := NewNacosLoader(nc, "DEFAULT_GROUP", "smart.gate.yaml", NewYamlDecoder())
	c, err := loader.Load()
	assert.Nil(t, err)
	t.Logf("%v", *c)
	err = loader.Watch(context.Background(), func(conf *Conf) {
		t.Logf("config change: %v", *conf)
	})
	assert.Nil(t, err)
	<-utility.WatchQuitSignal()
	t.Log("test finished")
}
