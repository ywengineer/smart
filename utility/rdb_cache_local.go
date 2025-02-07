package utility

import (
	"context"
	"github.com/go-gorm/caches/v4"
	lru "github.com/hashicorp/golang-lru/v2"
)

type memoryCacher struct {
	store *lru.Cache[string, any]
}

func (c *memoryCacher) size(s int) caches.Cacher {
	c.store, _ = lru.New[string, any](MaxInt(128, s))
	return c
}

func (c *memoryCacher) Get(ctx context.Context, key string, q *caches.Query[any]) (*caches.Query[any], error) {
	val, ok := c.store.Get(key)
	if !ok {
		return nil, nil
	}

	if err := q.Unmarshal(val.([]byte)); err != nil {
		return nil, err
	}

	return q, nil
}

func (c *memoryCacher) Store(ctx context.Context, key string, val *caches.Query[any]) error {
	res, err := val.Marshal()
	if err != nil {
		return err
	}
	c.store.Add(key, res)
	return nil
}

func (c *memoryCacher) Invalidate(ctx context.Context) error {
	c.store.Purge()
	return nil
}
