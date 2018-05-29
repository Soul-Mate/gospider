package core

import (
	"context"
	"sync"
	"github.com/Soul-Mate/gospider/http"
	"github.com/Soul-Mate/gospider/schedule"
	"github.com/Soul-Mate/gospider/spider"
	"github.com/Soul-Mate/gospider/download"
	"github.com/Soul-Mate/gospider/filter"
	"github.com/Soul-Mate/gospider/util/cmd"
	"github.com/Soul-Mate/gospider/util/conf"
	"github.com/Soul-Mate/gospider/util/log"
	"sort"
)

type Engine struct {
	wg        *sync.WaitGroup
	ctx       context.Context
	filter    filter.FilterInterface
	Spider    spider.SpiderInterface
	Spiders   map[string]spider.SpiderInterface
	Schedule  schedule.Scheduler
	ctxCancel context.CancelFunc
	download  *download.Downloader
}

func init() {
	cmdConf := cmd.InitCmd()
	conf.InitConfig(cmdConf.Path)
	log.InitLog()
	sort.Sort(conf.GlobalSharedConfig.Items)
}

func (e *Engine) Run() {
	e.StartSpider()

	e.download.Start()

	e.dispatch()
}

func (e *Engine) dispatch() {
	for {
		select {
		case <-e.ctx.Done():
			goto QUIT
		default:
			if e.Schedule.Len() <= 0 {
				goto QUIT
			}
			req := e.Schedule.Receive()
			e.download.Commit(req)
			resp := <-e.download.GetResponse()
			e.receiveSpiderGenerator(e.sendSpiderResponse(resp))
		}
	}
QUIT:
	e.exit()
}

func (e *Engine) sendRequest(r *http.Request) {
	// use filter
	if e.filter == nil {
		e.Schedule.Send(r)
		return
	} else {
		if e.filter.Contains(r.Options.Url) {
			return
		}
		e.filter.Add(r.Options.Url)
		e.Schedule.Send(r)
	}
}

func (e *Engine) exit() {
	e.filterDump()
	log.Printf("DEBUG", "spider quit\n")
}

func NewEngine(spider spider.SpiderInterface) *Engine {
	engine := new(Engine)
	ctx := context.Background()
	engine.wg = new(sync.WaitGroup)
	engine.ctx, engine.ctxCancel = context.WithCancel(ctx)
	engine.filter = newFilter()
	engine.Schedule = newSchedule()
	engine.Spider = spider
	engine.download = download.NewDownloader(engine.ctx, engine.wg)
	return engine
}

func newSchedule() schedule.Scheduler {
	return schedule.MakeSchedule(conf.GlobalSharedConfig.Schedule.Name)
}
