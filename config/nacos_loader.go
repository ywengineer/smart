package config

import (
	"context"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/pkg/errors"
	"log"
	"time"
)

type NacosLoader struct {
	Client  config_client.IConfigClient
	Group   string
	DataId  string
	Decoder Decoder
}

func (nl *NacosLoader) Load() (*Conf, error) {
	if err := nl.check(); err != nil {
		return nil, err
	}
	// get config
	content, err := nl.Client.GetConfig(vo.ConfigParam{Group: nl.Group, DataId: nl.DataId})
	if err != nil {
		return nil, errors.WithMessage(err, "load config content from nacos error")
	}
	conf := &Conf{}
	err = nl.Decoder.Unmarshal([]byte(content), conf)
	return conf, err
}

func (nl *NacosLoader) check() error {
	if nl.Client == nil {
		return errors.New("nacos client have not been initialized.")
	}
	if nl.Decoder == nil {
		return errors.New("nil config decoder is not allowed")
	}
	if len(nl.Group) == 0 || len(nl.DataId) == 0 {
		return errors.New("empty dataId and group is not allowed")
	}
	return nil
}

func (nl *NacosLoader) Watch(ctx context.Context, callback func(conf *Conf)) error {
	if err := nl.check(); err != nil {
		return err
	}
	p := vo.ConfigParam{
		DataId: nl.DataId,
		Group:  nl.Group,
		OnChange: func(namespace, group, dataId, data string) {
			conf := &Conf{}
			err := nl.Decoder.Unmarshal([]byte(data), conf)
			if err != nil {
				log.Printf("[NacosLoader] config changed. nacos loader parse error: %v\n", err)
			} else {
				callback(conf)
			}
		},
	}
	go func() {
		defer nl.Client.CancelListenConfig(p)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := ctx.Err(); err != nil {
					log.Printf("[NacosLoader] nacos loader watcher stopped. encounter an error: %v\n", err)
					return
				}
				time.Sleep(time.Second)
			}
		}
	}()
	return nl.Client.ListenConfig(p)
}
