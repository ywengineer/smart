package server_config

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
	"time"
)

type LocalLoader struct {
	Path    string `json:"path"`
	Decoder Decoder
}

func (ll *LocalLoader) Load() (*Conf, error) {
	if err := ll.check(); err != nil {
		return nil, err
	}
	if ll.Decoder == nil {
		fs := ll.Path[strings.LastIndex(ll.Path, ".")+1:]
		if strings.EqualFold(fs, "json") {
			ll.Decoder = &JSONDecoder{}
		} else if strings.EqualFold(fs, "yaml") || strings.EqualFold(fs, "yml") {
			ll.Decoder = &YamlDecoder{}
		} else {
			return nil, errors.Errorf("unsupported file : %s", ll.Path)
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

func (ll *LocalLoader) check() error {
	if len(ll.Path) == 0 {
		return errors.New("server_config file path is empty")
	}
	if !ll.isFileExist(ll.Path) {
		return errors.Errorf("server_config file[%s] is not exists", ll.Path)
	}
	return nil
}

func (ll *LocalLoader) Watch(ctx context.Context, callback func(conf *Conf)) error {
	if err := ll.check(); err != nil {
		return err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					conf, err := ll.Load()
					if err != nil || conf == nil {
						log.Printf("[LocalLoader] file changed. local loader parse error: %v\n", err)
					} else {
						callback(conf)
					}
				}
			case err = <-watcher.Errors:
				log.Printf("[LocalLoader] local loader watcher stopped. encounter an error: %v\n", err)
			case <-ctx.Done():
				return
			default:
				if err = ctx.Err(); err != nil {
					log.Printf("[LocalLoader] local loader watcher stopped. encounter an error: %v\n", err)
					return
				}
				time.Sleep(time.Second * 5)
			}
		}
	}()
	return watcher.Add(ll.Path)
}

func (ll *LocalLoader) isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
