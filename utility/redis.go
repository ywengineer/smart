package utility

import (
	"context"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"net/url"
)

func aa() {
	_, _ = redis.ParseURL("redis://127.0.0.1:6379/0")
	redis.NewClient(&redis.Options{})
	redis.NewClusterClient(&redis.ClusterOptions{})
	redis.NewUniversalClient(&redis.UniversalOptions{}).Get(context.Background(), "")
}

// NewRedis redisURL format redis://user:password@host:port/?db=0&node=host:port&node=host:port
//
// - when redis running in cluster mode, add query param node.
// - when redis running in simple mode, add query param db.
func NewRedis(redisURL string) redis.UniversalClient {
	u, err := url.Parse(redisURL)
	if err != nil {
		panic(errors.WithMessage(err, "parse redis url failed"))
		return nil
	}
	query := u.Query()
	opt := &redis.UniversalOptions{
		Addrs:        []string{u.Host},
		MinIdleConns: 2,
		DB:           0,
	}
	// multi nodes
	if len(query["node"]) > 0 {
		opt.Addrs = append(opt.Addrs, query["node"]...)
	}
	// select db
	opt.DB = MaxInt(0, QueryInt(query, "db"))
	// use auth
	if u.User != nil {
		opt.Username = u.User.Username()
		opt.Password, _ = u.User.Password()
	}
	//
	return redis.NewUniversalClient(opt)
}
