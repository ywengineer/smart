package server_config

import (
	"github.com/bytedance/sonic"
	"gopkg.in/yaml.v3"
)

type Conf struct {
	Network           string                 `json:"network" yaml:"network"`
	Address           string                 `json:"address" yaml:"address"`
	Workers           int                    `json:"workers" yaml:"workers"`
	WorkerLoadBalance string                 `json:"load_balance" yaml:"load-balance"`
	Attach            map[string]interface{} `json:"attach" yaml:"attach"`
}

type Decoder interface {
	Unmarshal(buf []byte, val interface{}) error
}

type JSONDecoder struct {
}

func (d *JSONDecoder) Unmarshal(buf []byte, val interface{}) error {
	return sonic.Unmarshal(buf, val)
}

type YamlDecoder struct {
}

func (d *YamlDecoder) Unmarshal(buf []byte, val interface{}) error {
	return yaml.Unmarshal(buf, val)
}
