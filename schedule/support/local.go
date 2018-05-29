package support

import (
	"github.com/Soul-Mate/gospider/http"
)

type SingleSchedule struct {
	RequestQueue []*http.Request
}

func (s *SingleSchedule) Send(r *http.Request) {
	s.RequestQueue = append(s.RequestQueue, r)
}

func (s *SingleSchedule) Receive() *http.Request {
	ret := s.RequestQueue[0]
	s.RequestQueue = s.RequestQueue[1:]
	return ret
}

func (s *SingleSchedule) Len() int  {
	return len(s.RequestQueue)
}

func (s *SingleSchedule) Close()  {

}