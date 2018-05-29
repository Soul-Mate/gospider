package support

import (
	"github.com/Soul-Mate/gospider/http"
	"github.com/Soul-Mate/gospider/queue/request"
)

type RedisSchedule struct {
	RequestQueue   *request.RedisRequestQueue
}

func (rs *RedisSchedule) Send(r *http.Request) {
	rs.RequestQueue.Push(r)
}

func (rs *RedisSchedule) Receive() *http.Request  {
	return rs.RequestQueue.Pop()
}

func (rs *RedisSchedule) Len() int {
	return rs.RequestQueue.Len()
}

func (rs *RedisSchedule) Close() {
	rs.RequestQueue.Close()
}
