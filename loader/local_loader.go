package loader

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
	"time"
)

type localLoader struct {
	path    string
	decoder Decoder
}

func NewLocalLoader(path string) SmartLoader {
	return &localLoader{path: path}
}

func (ll *localLoader) Load(out interface{}) error {
	if err := ll.check(); err != nil {
		return err
	}
	if ll.decoder == nil {
		fs := ll.path[strings.LastIndex(ll.path, ".")+1:]
		if strings.EqualFold(fs, "json") {
			ll.decoder = &jsonDecoder{}
		} else if strings.EqualFold(fs, "yaml") || strings.EqualFold(fs, "yml") {
			ll.decoder = &yamlDecoder{}
		} else {
			return errors.Errorf("unsupported file : %s", ll.path)
		}
	}
	data, err := os.ReadFile(ll.path)
	if err != nil {
		return err
	}
	return ll.decoder.Unmarshal(data, out)
}

func (ll *localLoader) check() error {
	if len(ll.path) == 0 {
		return errors.New("loader file path is empty")
	}
	if !ll.isFileExist(ll.path) {
		return errors.Errorf("loader file[%s] is not exists", ll.path)
	}
	return nil
}

func (ll *localLoader) Watch(ctx context.Context, callback WatchCallback) error {
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
					conf := &Conf{}
					if ll.Load(conf) != nil {
						log.Printf("[localLoader] file changed. local loader parse error: %v\n", err)
					} else {
						_ = callback(conf)
					}
				}
			case err = <-watcher.Errors:
				log.Printf("[localLoader] local loader watcher stopped. encounter an error: %v\n", err)
			case <-ctx.Done():
				return
			default:
				if err = ctx.Err(); err != nil {
					log.Printf("[localLoader] local loader watcher stopped. encounter an error: %v\n", err)
					return
				}
				time.Sleep(time.Second * 5)
			}
		}
	}()
	return watcher.Add(ll.path)
}

func (ll *localLoader) isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
