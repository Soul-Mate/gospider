package core

import (
	"github.com/Soul-Mate/gospider/filter"
	"github.com/Soul-Mate/gospider/util/conf"
	"strings"
	"strconv"
)

func newFilter() filter.FilterInterface {
	c := conf.GlobalSharedConfig.Filter
	if !c.Usage {
		return nil
	}
	name := strings.ToLower(c.Name)
	switch name {
	case "local": // will dump local file
		return filter.MakeFilter(parseFilterLength(c.Length), "local")
	case "redis": // use redis bitset
		return filter.MakeFilter(0, "redis")
	default:
		return nil
	}
}

// dump the filter file, only support local
func (e *Engine) filterDump() {
	fConf := conf.GlobalSharedConfig.Filter
	if e.filter != nil && fConf.Dump && fConf.DumpPath != "" && strings.ToLower(fConf.Name) == "local" {
		e.filter.(*filter.LocalFilter).Dump(fConf.DumpPath)
	}
}

// parse dump file size
// the length support: int value and string value
// example: 1024  or 1kb, 1024 * 1024 or 1gb
func parseFilterLength(length interface{}) uint {
	var (
		ret           uint
		defaultLength = 8*1024*1024*500 - 1 // 500MB
	)
	switch length.(type) {
	case int:
		if (length.(int)) <= 0 {
			ret = uint(defaultLength)
		} else {
			ret = uint(length.(int))
		}
	case string:
		str := length.(string)
		strlen := len(str)
		neg := str[strlen-2:]
		str = str[:strlen-2]
		if n, err := strconv.Atoi(str); err != nil {
			ret = uint(defaultLength)
		} else {
			switch neg {
			case "kb":
				ret = uint(8 * 1024 * n)
			case "KB":
				ret = uint(8 * 1024 * n)
			case "mb":
				ret = uint(8 * 1024 * 1024 * n)
			case "MB":
				ret = uint(8 * 1024 * 1024 * n)
			case "gb":
				ret = uint(8 * 1024 * 1024 * 1024 * n)
			case "GB":
				ret = uint(8 * 1024 * 1024 * 1024 * n)
			default:
				ret = uint(defaultLength)
			}
		}
	default:
		ret = uint(defaultLength)
	}
	return ret
}
