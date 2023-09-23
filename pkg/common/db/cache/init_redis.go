package cache

import (
	"errors"
	"github.com/chdlvy/go-chatFrame/pkg/common/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis() (redis.UniversalClient, error) {
	var rdb redis.UniversalClient
	if len(config.Config.Redis.Address) == 0 {
		return nil, errors.New("redis address is empty")
	}
	if len(config.Config.Redis.Address) > 1 {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:      config.Config.Redis.Address,
			Username:   config.Config.Redis.Username,
			Password:   config.Config.Redis.Password,
			PoolSize:   50,
			MaxRetries: 10,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:       config.Config.Redis.Address[0],
			Username:   config.Config.Redis.Username,
			Password:   config.Config.Redis.Password,
			DB:         0,
			PoolSize:   100,
			MaxRetries: 10,
		})
	}
	return rdb, nil
}
