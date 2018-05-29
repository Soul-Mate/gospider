package filter

import (
	"github.com/willf/bloom"
	"os"
	"github.com/Soul-Mate/gospider/util/conf"
)

type LocalFilter struct {
	Name        string
	BloomFilter *bloom.BloomFilter
}

func (lf *LocalFilter) Add(url string) {
	lf.BloomFilter.AddString(url)
}

func (lf *LocalFilter) AddByte(url []byte) {
	lf.BloomFilter.Add(url)
}

func (lf *LocalFilter) Contains(url string) bool {
	return lf.BloomFilter.TestString(url)
}

func (lf *LocalFilter) ContainsByte(url []byte) bool {
	return lf.BloomFilter.Test(url)
}

func (lf *LocalFilter) Dump(filepath string) error {
	var (
		err error
		fd  *os.File
	)
	_, err = os.Stat(filepath)
	if err == nil {
		if fd, err = os.Open(filepath); err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		if fd, err = os.Create(filepath); err != nil {
			return err
		}
	} else {
		return err
	}

	defer fd.Close()
	_, err = lf.BloomFilter.WriteTo(fd)
	return err
}


func NewLocalFilter(length uint) *LocalFilter {
	lf := new(LocalFilter)
	lf.Name = "local"
	lf.BloomFilter = bloom.New(length, 5)
	lf.fromDumpFile()
	return lf
}

func (lf *LocalFilter) fromDumpFile() {
	if conf.GlobalSharedConfig.Filter.Dump && conf.GlobalSharedConfig.Filter.DumpPath != "" {
		if fd, err := os.Open(conf.GlobalSharedConfig.Filter.DumpPath); err != nil {
			return
		} else {
			defer fd.Close()
			lf.BloomFilter.ReadFrom(fd)
			return
		}
	}
	return
}
