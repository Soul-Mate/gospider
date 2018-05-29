package core

import (
	"github.com/Soul-Mate/gospider/spider"
	"github.com/Soul-Mate/gospider/http"
)

func (e *Engine) StartSpider() {
	for rs := range e.Spider.StartRequests() {
		for _, req := range rs {
			e.sendRequest(req)
		}
	}
}

// send response to spider
func (e *Engine) sendSpiderResponse(r *http.Response) (<-chan interface{}) {
	var (
		callback spider.CallBackFunc
	)
	if r == nil {
		return nil
	}

	if callback = e.Spider.FindCallBack(r.CallBackName); callback == nil {
		return nil
	}

	return callback(e.Spider, r)
}

// receive spider generator requests and items
func (e *Engine) receiveSpiderGenerator(generator <-chan interface{}) {
	for v := range generator {
		switch v.(type) {
		case *http.Request:
			e.receiveRequest(v.(*http.Request))
		case []*http.Request:
			e.receiveRequest(v.([]*http.Request)...)
		case spider.Item:
			e.receiveItem(e.Spider, v.(spider.Item))
		case []spider.Item:
			e.receiveItem(e.Spider, v.([]spider.Item)...)
		case *spider.Item:
			e.receivePointerItem(e.Spider, v.(*spider.Item))
		case []*spider.Item:
			e.receivePointerItem(e.Spider, v.([]*spider.Item)...)
		default:
		}
	}
}

// receive request.
// these requests will be send to request queue
func (e *Engine) receiveRequest(rs ...*http.Request) {
	for _, r := range rs {
		e.sendRequest(r)
	}
}

// receive item.
// these items will be handler over to the pipeline for processing.
func (e *Engine) receiveItem(sp spider.SpiderInterface, items ...spider.Item) {
	for _, it := range items {
		e.dispatchPipeline(&it, sp)
	}
}

// receive item pointer.
func (e *Engine) receivePointerItem(sp spider.SpiderInterface, items ...*spider.Item) {
	for _, it := range items {
		e.dispatchPipeline(it, sp)
	}
}
