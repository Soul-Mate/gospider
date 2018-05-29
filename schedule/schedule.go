package schedule

import (
	"github.com/Soul-Mate/gospider/http"
	"github.com/Soul-Mate/gospider/schedule/support"
	"github.com/Soul-Mate/gospider/queue/request"
)

type Scheduler interface {
	Send(r *http.Request)
	Receive() *http.Request
	Close()
	Len() int
}

// make schedule by name
func MakeSchedule(name string) Scheduler {
	switch name {
	case "local":
		return NewLocalSchedule()
	case "redis":
		return NewRedisSchedule()
	default:
		return nil
	}
}

func NewRedisSchedule() *support.RedisSchedule {
	rs := new(support.RedisSchedule)
	rs.RequestQueue = request.NewRedisRequestQueue()
	return rs
}

func NewLocalSchedule() *support.SingleSchedule {
	s := &support.SingleSchedule{
		RequestQueue: make([]*http.Request, 0),
	}
	return s
}
