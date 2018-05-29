package download

import (
	"github.com/Soul-Mate/gospider/http"
	"context"
	"sync"
	"github.com/Soul-Mate/gospider/util/conf"
	"time"
	"github.com/Soul-Mate/gospider/util/log"
)

type Downloader struct {
	wg   *sync.WaitGroup
	ctx  context.Context
	in   chan *http.Request
	out  chan *http.Response
	tick <-chan time.Time
}

func NewDownloader(ctx context.Context, wg *sync.WaitGroup) *Downloader {
	downloader := new(Downloader)
	downloader.wg = wg
	downloader.ctx = ctx
	downloader.in = make(chan *http.Request)
	downloader.out = make(chan *http.Response, 10)
	downloader.setDownloadDelay()
	return downloader
}

// set download delay time
// support format: 1, 0.1 int or float value unit is second,
// and to support 0.1s 1s string value
func (d *Downloader) setDownloadDelay() {
	var (
		err          error
		delay        interface{}
		tickDuration time.Duration
	)

	delay = conf.GlobalSharedConfig.Downloader.Delay

	switch delay.(type) {
	case int:
		if delay.(int) == 0 {
			d.tick = nil
			return
		}
		tickDuration = time.Second * time.Duration(delay.(int))
	case float64:
		if delay.(float64) <= 0 {
			d.tick = nil
			return
		}
		tickDuration = time.Duration(delay.(float64) * float64(time.Second))
	case string:
		if tickDuration, err = time.ParseDuration(delay.(string)); err != nil {
			tickDuration = time.Second / 3
		}
	default:
		tickDuration = time.Second / 3
	}

	d.tick = time.Tick(tickDuration)
}

func (d *Downloader) Commit(r *http.Request) {
	go func() {
		d.in <- r
	}()
}

func (d *Downloader) GetResponse() <-chan *http.Response {
	return d.out
}

func (d *Downloader) Start() {
	for i := 0; i < conf.GlobalSharedConfig.Downloader.Concurrent; i++ {
		d.wg.Add(1)
		go d.workerStart()
	}
}

// start user define worker
// multiple worker will grab the request submitted by the engine.
// if fetch ok. response will be send to output channel
func (d *Downloader) workerStart() {
	var exit bool
	for {
		select {
		case req := <-d.in:
			if req != nil {
				if resp, err := d.fetch(req); err != nil {
					log.Printf("DEBUG", "download: %s, error: %s\n", req.Options.Url, err.Error())
					d.out <- nil
				} else {
					d.out <- resp
				}
			}
		case <-d.ctx.Done():
			exit = true
			break
		}
		if exit {
			break
		}
	}
	log.Printf("DEBUG", "downloader exit\n")
	d.wg.Done()
}

// exec request. support delay download
func (d *Downloader) fetch(req *http.Request) (response *http.Response, err error) {
	if d.tick != nil {
		<-d.tick
	}
	log.Printf("DEBUG", "fetch: %s\n", req.Options.Url)
	return req.Do()
}
