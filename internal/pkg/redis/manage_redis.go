package redis

import (
	"github.com/garyburd/redigo/redis"
)

type ManageRedis struct {
}

const (
	HOST = "127.0.0.1:6379"
	// userid:phone的映射
	IDPHONE = "phone:uid"
)

var pool *redis.Pool

func GetredisConn() redis.Conn {

	if pool == nil {
		pool = &redis.Pool{
			MaxIdle:     16,
			MaxActive:   0,
			IdleTimeout: 300,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", HOST)
			},
		}
	}
	// 从连接池中拿一个连接出去
	return pool.Get()
}

// 可能要重写逻辑
func InitRedis() {
	rdb := GetredisConn()
	defer rdb.Close()
	rdb.Send("multi")
	rdb.Send("flushall")
	rdb.Send("set", "user:max_id", "0")
	rdb.Send("select", "1")
	rdb.Send("set", "video:max_id", "0")
	rdb.Send("set", "comment:max_id", "0")
	rdb.Send("exec")

}
