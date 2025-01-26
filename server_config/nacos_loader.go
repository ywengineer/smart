package server_config

import (
	"context"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/pkg/errors"
	"github.com/ywengineer/smart/utility"
	"log"
	"time"
)

type nacosLoader struct {
	nc      config_client.IConfigClient
	group   string
	dataId  string
	decoder Decoder
}

func NewNacosLoader(nacos config_client.IConfigClient, group string, dataId string, decoder Decoder) SmartLoader {
	return &nacosLoader{
		nc:      nacos,
		dataId:  utility.IfEmptyStr(dataId, "smart.server.yaml"),
		group:   utility.IfEmptyStr(group, "DEFAULT_GROUP"),
		decoder: decoder,
	}
}

func NewDefaultNacosLoader(nacos config_client.IConfigClient, dataId string, decoder Decoder) SmartLoader {
	return NewNacosLoader(nacos, "DEFAULT_GROUP", dataId, decoder)
}

// NewNacosClient
// contextPath, nacos server context path
// the logLevel must be debug,info,warn,error, default value is info
func NewNacosClient(ipAddr string, port uint64, contextPath string,
	timeoutMs uint64,
	namespace, user, password, logLevel string,
) (config_client.IConfigClient, error) {
	// create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ipAddr, port, constant.WithContextPath(contextPath)),
	}
	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(namespace),
		constant.WithTimeoutMs(timeoutMs),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithUsername(user),
		constant.WithPassword(password),
		constant.WithLogLevel(logLevel),
	)
	// create server_config client
	return clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
}

func (nl *nacosLoader) Load(out interface{}) error {
	if err := nl.check(); err != nil {
		return err
	}
	// get server_config
	content, err := nl.nc.GetConfig(vo.ConfigParam{Group: nl.group, DataId: nl.dataId})
	if err != nil {
		return errors.WithMessage(err, "load server_config content from nacos error")
	}
	return nl.decoder.Unmarshal([]byte(content), out)
}

func (nl *nacosLoader) check() error {
	if nl.nc == nil {
		return errors.New("nacos client have not been initialized.")
	}
	if nl.decoder == nil {
		return errors.New("nil server_config decoder is not allowed")
	}
	if len(nl.group) == 0 || len(nl.dataId) == 0 {
		return errors.New("empty dataId and group is not allowed")
	}
	return nil
}

func (nl *nacosLoader) Watch(ctx context.Context, callback WatchCallback) error {
	if err := nl.check(); err != nil {
		return err
	}
	p := vo.ConfigParam{
		DataId: nl.dataId,
		Group:  nl.group,
		OnChange: func(namespace, group, dataId, data string) {
			conf := &Conf{}
			err := nl.decoder.Unmarshal([]byte(data), conf)
			if err != nil {
				log.Printf("[nacosLoader] server_config changed. nacos loader parse error: %v\n", err)
			} else {
				_ = callback(conf)
			}
		},
	}
	go func() {
		defer func() {
			_ = nl.nc.CancelListenConfig(p)
			nl.nc.CloseClient()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := ctx.Err(); err != nil {
					log.Printf("[nacosLoader] nacos loader watcher stopped. encounter an error: %v\n", err)
					return
				}
				time.Sleep(time.Second)
			}
		}
	}()
	return nl.nc.ListenConfig(p)
}
