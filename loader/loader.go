package loader

import (
	"context"
	"reflect"
)

type WatchCallback func(data string) error

type SmartLoader interface {
	Load(outPointer interface{}) error
	Watch(ctx context.Context, callback WatchCallback) error
	Unmarshal(data []byte, out interface{}) error
}

func NewValueLoader(value interface{}) SmartLoader {
	return &valueLoader{value: value}
}

type valueLoader struct {
	value interface{}
}

func (vl *valueLoader) Unmarshal(data []byte, out interface{}) error {
	return nil
}

func (vl *valueLoader) Load(outPointer interface{}) error {
	reflect.ValueOf(outPointer).Elem().Set(reflect.ValueOf(vl.value).Elem())
	//reflect.ValueOf(outPointer).Set(reflect.ValueOf(vl.value))
	return nil
}

func (vl *valueLoader) Watch(ctx context.Context, callback WatchCallback) error {
	return nil
}
