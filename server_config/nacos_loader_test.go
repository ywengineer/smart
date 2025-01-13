package server_config

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNacosLoader(t *testing.T) {
	nc := newNacosTestClient(t)
	loader := NewNacosLoader(nc, "DEFAULT_GROUP", "smart.server.yaml", NewYamlDecoder())
	c, err := loader.Load()
	assert.Nil(t, err)
	t.Logf("%v", *c)
}

func newNacosTestClient(t *testing.T) config_client.IConfigClient {
	//create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("192.168.0.15", 8848, constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithUsername("nacos"),
		constant.WithPassword("vspn"),
		constant.WithLogLevel("debug"),
	)
	// create server_config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	assert.Nil(t, err)
	return client
}
