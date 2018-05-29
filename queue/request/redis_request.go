package request

import (
	"github.com/gomodule/redigo/redis"
	"github.com/Soul-Mate/gospider/http"
	"github.com/Soul-Mate/gospider/util"
)

const DEFAULT_KEY = "wait_crawl_queue"

type RedisRequestQueue struct {
	pool *redis.Pool
}

func (rq *RedisRequestQueue) Push(req *http.Request) {
	var (
		err     error
		content []byte
	)
	conn := rq.pool.Get()
	defer conn.Close()

	// marshal request
	if content, err = req.Marshal(); err == nil {
		conn.Do("LPUSH", DEFAULT_KEY, string(content))
	}
}

func (rq *RedisRequestQueue) Pop() *http.Request {
	var (
		ok      bool
		err     error
		content []byte
		reply   interface{}
	)
	conn := rq.pool.Get()
	defer conn.Close()

	if reply, err = conn.Do("LPOP", DEFAULT_KEY); err != nil {
		return nil
	}
	if content, ok = reply.([]byte); !ok {
		return nil
	}
	req := new(http.Request)
	if err = req.Unmarshal(content); err != nil {
		return nil
	}
	return req
}

func (rq *RedisRequestQueue) Len() int {
	conn := rq.pool.Get()
	defer conn.Close()
	if reply, err := conn.Do("LLEN", DEFAULT_KEY); err != nil {
		return 0
	} else {
		return int(reply.(int64))
	}
}

func (rq *RedisRequestQueue) Close() {
	rq.pool.Close()
}

func NewRedisRequestQueue() *RedisRequestQueue {
	rd := new(RedisRequestQueue)
	rd.pool = util.GetSharedPool()
	return rd
}
