package filter

import (
	"testing"
	"github.com/willf/bloom"
	"os"
	"bufio"
	"io"
	"errors"
)

func TestBloomFilter(t *testing.T) {
	f := bloom.New(1*1024*1024, 5)
	f.AddString("love")
	if !f.TestString("love") {
		t.Error("TestString error")
	}
}

func TestBloomFilterForFile(t *testing.T) {
	b := bloom.New(1*1024*1024, 5)
	err := addToBetSet(b)
	if err != nil {
		t.Error(err.Error())
	}
	err = testBetSet(b)
	if err != nil {
		t.Error(err.Error())
	}

}

func TestBloomFilterAop(t *testing.T) {
	// 1 GB
	b := bloom.New(8*1024*1024*1024, 5)
	err := addToBetSet(b)
	if err != nil {
		t.Error(err.Error())
	}
	betSetAop(b)
}

func TestBloomFilterForAop(t *testing.T) {
	// 1 GB
	b := bloom.New(8*1024*1024*1024, 5)
	err := addToBetSet(b)
	if err != nil {
		t.Error(err.Error())
	}
	betSetAop(b)

	b, err = getBloomForAopFile()
	if err != nil {
		t.Error(err.Error())
	}
	testBetSet(b)
}

func getStream() (*bufio.Reader, error) {
	f, err := os.Open("test.txt")
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	return buf, nil
}

func addToBetSet(b *bloom.BloomFilter) error {
	stream, err := getStream()
	if err != nil {
		return err
	}
	for {
		line, _, err := stream.ReadLine()
		b.Add(line)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
	}
	return nil
}

func testBetSet(b *bloom.BloomFilter) error {
	stream, err := getStream()
	if err != nil {
		return err
	}
	for {
		line, _, err := stream.ReadLine()
		if !b.Test(line) {
			return errors.New("error")
		}
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
	}
	return nil
}

func betSetAop(b *bloom.BloomFilter) error {
	f, err := os.Create("bloomfilter.data")
	if err != nil {
		return err
	}
	_, err = b.WriteTo(f)
	return err
}

func getBloomForAopFile() (*bloom.BloomFilter, error) {
	f, err := os.Open("bloomfilter.data")
	if err != nil {
		return nil, err
	}
	b := &bloom.BloomFilter{}

	_, err = b.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}
