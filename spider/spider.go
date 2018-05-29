package spider

import (
	"github.com/Soul-Mate/gospider/http"
	"sync"
)

type SpiderInterface interface {
	StartRequests() chan []*http.Request
	Generator(data ...interface{}) chan interface{}
	FindCallBack(name string) CallBackFunc
	SearchItem(name string) ItemPipeline
}

type CallBackFunc func(spider SpiderInterface, response *http.Response) chan interface{}

type Spider struct {
	wg                *sync.WaitGroup
	Name              string
	AllowedDomains    []string
	StartUrls         []string
	Quit              chan interface{}
	CallBackFunctions map[string]CallBackFunc
	Items             map[string]ItemPipeline
}

func NewSpider(spiderName string, startUrls []string, funcs map[string]CallBackFunc, items map[string]ItemPipeline) *Spider {
	spider := new(Spider)
	spider.wg = new(sync.WaitGroup)
	spider.Name = spiderName
	spider.StartUrls = startUrls
	spider.CallBackFunctions = funcs
	spider.Items = items
	return spider
}

func (s *Spider) StartRequests() chan []*http.Request {
	c := make(chan []*http.Request, len(s.StartUrls))
	requests := make([]*http.Request, len(s.StartUrls)-1)
	go func() {
		for _, url := range s.StartUrls {
			req := http.NewRequest("GET", url, nil)
			req.SpiderName = s.Name
			req.CallBackName = "Response"
			requests = append(requests, req)
		}
		c <- requests
		close(c)
	}()
	return c
}

func (s *Spider) Generator(data ...interface{}) chan interface{} {
	ret := make(chan interface{}, len(data))
	go func() {
		for _, v := range data {
			ret <- v
		}
		close(ret)
	}()
	return ret
}

func (s *Spider) FindCallBack(name string) CallBackFunc {
	if callback, ok := s.CallBackFunctions[name]; !ok {
		return nil
	} else {
		return callback
	}
}

func (s *Spider) SearchItem(name string) ItemPipeline {
	if itemer, ok := s.Items[name]; !ok {
		return nil
	} else {
		return itemer
	}
}
