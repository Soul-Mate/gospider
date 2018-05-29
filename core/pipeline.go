package core

import (
	"github.com/Soul-Mate/gospider/util/conf"
	"github.com/Soul-Mate/gospider/spider"
)

func (e *Engine) dispatchPipeline(item *spider.Item, s spider.SpiderInterface) {
	items := conf.GlobalSharedConfig.Items
	N := len(items)
	for i := 0; i < N; i++ {
		if it := s.SearchItem(items[i].Name); it != nil {
			if item = it.ProcessItem(item, s); item.IsSkip() {
				break
			}
		}
	}
}