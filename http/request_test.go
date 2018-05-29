package http

import (
	"testing"
)

func TestNewRequest(t *testing.T) {
	opt := &Options{
		Json: true,
		Body: map[string]string{},
	}
	r := NewRequest("get", "http://www.baidu.com", opt)
	resp, _ := r.Do()
	println(resp.StatusCode)
}
