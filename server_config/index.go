package server_config

import (
	"github.com/bytedance/sonic"
	"gopkg.in/yaml.v3"
)

type Conf struct {
	ServiceName       string                 `json:"service_name" yaml:"service-name"`
	Weight            int                    `json:"weight" yaml:"weight"`
	Network           string                 `json:"network" yaml:"network"`
	Address           string                 `json:"address" yaml:"address"`
	Workers           int                    `json:"workers" yaml:"workers"`
	WorkerLoadBalance string                 `json:"load_balance" yaml:"load-balance"`
	Metadata          map[string]interface{} `json:"metadata" yaml:"metadata"`
}

type Decoder interface {
	Unmarshal(buf []byte, val interface{}) error
}

func NewJsonDecoder() Decoder {
	return &jsonDecoder{}
}

func NewYamlDecoder() Decoder {
	return &yamlDecoder{}
}

type jsonDecoder struct {
}

func (d *jsonDecoder) Unmarshal(buf []byte, val interface{}) error {
	return sonic.Unmarshal(buf, val)
}

type yamlDecoder struct {
}

func (d *yamlDecoder) Unmarshal(buf []byte, val interface{}) error {
	return yaml.Unmarshal(buf, val)
}
