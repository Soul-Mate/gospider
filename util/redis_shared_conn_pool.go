package util

import (
	"github.com/gomodule/redigo/redis"
	"github.com/Soul-Mate/gospider/util/conf"
)

var sharedPool *redis.Pool

func GetSharedPool() *redis.Pool {
	if sharedPool == nil {
		sharedPool = newPool()
	}
	return sharedPool
}

func newPool() *redis.Pool {
	redisPool := &redis.Pool{}
	redisPool.MaxIdle = conf.GlobalSharedConfig.Connection.Redis.ConnectionNums
	redisPool.Dial = func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", conf.GlobalSharedConfig.Connection.Redis.Addr)
		if err != nil {
			return nil, err
		}
		redis.DialDatabase(conf.GlobalSharedConfig.Connection.Redis.DB)
		return c, nil
	}
	return redisPool
}
