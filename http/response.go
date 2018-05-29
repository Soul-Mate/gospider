package http

import (
	"net/http"
	"encoding/json"
)

type Response struct {
	Ok           bool
	StatusCode   int
	SpiderName   string // spider name
	CallBackName string // response after call func name
}

func WrapResponse(req *Request, resp *http.Response) (*Response) {
	res := new(Response)
	res.SpiderName = req.SpiderName
	res.CallBackName = req.CallBackName
	return res
}

func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Response) UnMarshal(bs []byte) error{
	return json.Unmarshal(bs, r)
}
