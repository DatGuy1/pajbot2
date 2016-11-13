package models

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/common/config"
)

// RedisManager keeps the pool of redis connections
type RedisManager struct {
	Pool *redis.Pool
}

// InitRedis connects to redis and returns redis client
func InitRedis(config *config.Config) {
	r := &RedisManager{}
	pool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", config.RedisHost)
		if err != nil {
			return nil, err
		}
		if config.RedisDatabase >= 0 {
			_, err = c.Do("SELECT", config.RedisDatabase)
			if err != nil {
				return nil, err
			}
		}
		return c, err
	}, 69)
	// forsenGASM
	r.Pool = pool
	db = r
}

// GetRedisConn XD
// WARNING: You're responsible for closing the connection
func GetRedisConn() redis.Conn {
	return db.Pool.Get()
}

// DB xD
var db *RedisManager
