package spider

type ItemPipeline interface {
	Name() string // get pipeline name
	ProcessItem(it *Item, sp SpiderInterface) *Item
	CloseSpider(SpiderInterface)
}

type Field interface{}

type Item struct {
	Data map[string]interface{}
	skip bool
}

func (it *Item) SkipItem() {
	it.skip = true
}

func (it *Item) IsSkip() bool {
	if it.skip {
		return true
	}
	return false
}

func NewItem() *Item {
	it := new(Item)
	it.Data = make(map[string]interface{})
	it.skip = false
	return it
}
