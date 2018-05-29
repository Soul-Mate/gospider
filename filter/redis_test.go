package filter

import (
	"testing"
)

func TestNewRedisFilter(t *testing.T) {
	rf := NewRedisFilter()
	if rf.Pool == nil {
		t.Error("pool nil")
	}
	if rf.offset != MAX_SIZE {
		t.Error("offset error")
	}
	if rf.hashNum != 5 {
		t.Error("hash num error")
	}
}

func TestRedisFilter(t *testing.T) {
	rf := NewRedisFilter()
	rf.Add("http://www.example.com")
	if !rf.Contains("http://www.example.com") {
		t.Error("filter error")
	}

	if rf.Contains("https://www.example.com") {
		t.Error("filter error")
	}
}