package config

import (
	"context"
	"github.com/pkg/errors"
)

type Loader interface {
	Load() (*Conf, error)
	Watch(ctx context.Context, callback func(conf *Conf)) error
}

type ValueLoader struct {
	Conf *Conf
}

func (vl *ValueLoader) Load() (*Conf, error) {
	if vl.Conf == nil {
		return nil, errors.New("conf value is nil")
	}
	return vl.Conf, nil
}

func (vl *ValueLoader) Watch(ctx context.Context, callback func(conf *Conf)) error {
	callback(vl.Conf)
	return nil
}
