package config

import (
	"github.com/pkg/errors"
	"os"
	"strings"
)

type LocalLoader struct {
	Path    string `json:"path"`
	Decoder Decoder
}

func (ll *LocalLoader) Load() (*Conf, error) {
	if len(ll.Path) == 0 {
		return nil, errors.New("config file path is empty")
	}
	if !ll.isFileExist(ll.Path) {
		return nil, errors.Errorf("config file[%s] is not exists", ll.Path)
	}
	if ll.Decoder == nil {
		fs := ll.Path[strings.LastIndex(ll.Path, ".")+1:]
		if strings.EqualFold(fs, "json") {
			ll.Decoder = &JSONDecoder{}
		} else if strings.EqualFold(fs, "yaml") || strings.EqualFold(fs, "yml") {
			ll.Decoder = &YamlDecoder{}
		} else {
			return nil, errors.Errorf("not supported file : %s", ll.Path)
		}
	}
	data, err := os.ReadFile(ll.Path)
	if err != nil {
		return nil, err
	}
	conf := &Conf{}
	err = ll.Decoder.Unmarshal(data, conf)
	return conf, err
}

func (ll *LocalLoader) isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
