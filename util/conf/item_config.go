package conf

type ItemPipelineConfig struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

type Items []ItemPipelineConfig

func (i Items) Len() int {
	return len(i)
}

func (i Items) Swap(x, y int) {
	i[x], i[y] = i[y], i[x]
}

func (i Items) Less(x, y int) bool {
	return i[x].Priority < i[y].Priority
}
