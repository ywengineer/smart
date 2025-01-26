package server_config

import (
	"context"
)

type WatchCallback func(c interface{}) error

type SmartLoader interface {
	Load(outPointer interface{}) error
	Watch(ctx context.Context, callback WatchCallback) error
}

type ValueLoader struct {
	Conf *Conf
}

func (vl *ValueLoader) Load(outPointer interface{}) error {
	outPointer = vl.Conf
	return nil
}

func (vl *ValueLoader) Watch(ctx context.Context, callback WatchCallback) error {
	return callback(vl.Conf)
}
