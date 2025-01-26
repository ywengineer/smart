package loader

import (
	"context"
)

type WatchCallback func(c interface{}) error

type SmartLoader interface {
	Load(outPointer interface{}) error
	Watch(ctx context.Context, callback WatchCallback) error
}

func NewValueLoader(value interface{}) SmartLoader {
	return &valueLoader{value: value}
}

type valueLoader struct {
	value interface{}
}

func (vl *valueLoader) Load(outPointer interface{}) error {
	outPointer = vl.value
	return nil
}

func (vl *valueLoader) Watch(ctx context.Context, callback WatchCallback) error {
	return callback(vl.value)
}
