package http

import (
	"net/http"
	"strings"
	"bytes"
	"net/url"
	"net/http/cookiejar"
	"errors"
	"encoding/json"
	"io"
	"time"
	"github.com/Soul-Mate/gospider/util/conf"
)

type Options struct {
	Url       string // request url
	Method    string // request method
	Json      interface{}
	Data      interface{}
	XML       interface{}
	Headers   map[string]string
	Cookies   []*http.Cookie    // cookie
	Body      map[string]string // request body ["key" => "value"]
	ProxyAddr string            // request proxy
}

type Request struct {
	Options      *Options
	SpiderName   string      // spider name
	CallBackName string      // response after call func name
	Header       http.Header // request header
	Seen         string      // request seen
}

// make request
func NewRequest(method string, url string, opt *Options) (*Request) {
	if opt == nil {
		opt = new(Options)
	}
	opt.Url = url
	opt.Method = strings.ToUpper(method)
	req := new(Request)
	req.Options = opt
	return req
}

// make go request
func (r *Request) newGoRequest() (req *http.Request, err error) {
	if r.Options.Method == http.MethodGet {
		req, err = r.buildQueryStringRequest()
	}
	if r.Options.Method == http.MethodPost {
		if r.Options.Json != nil {
			req, err = r.buildJsonRequest()
		}
		if r.Options.Data != nil {
			req, err = r.buildDataRequest()
		}
		if r.Options.XML != nil {
			// TODO implement me
			req, err = nil, errors.New("no implement XML")
		}
	}

	// set header
	if err == nil {
		r.setMIMEHeader(req)
	}
	return
}

// Content-Type:application/json
func (r *Request) buildJsonRequest() (req *http.Request, err error) {
	var body io.Reader
	if bs, err := json.Marshal(r.Options.Body); err != nil {
		body = nil
	} else {
		body = bytes.NewReader(bs)
	}
	if req, err = http.NewRequest(r.Options.Method, r.Options.Url, body); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	return
}

// Content-Type:application/x-www-form-urlencoded
func (r *Request) buildDataRequest() (req *http.Request, err error) {
	var body io.Reader
	data := encodeBodyValues(r.Options.Body)
	body = strings.NewReader(data.Encode())
	if req, err = http.NewRequest(r.Options.Method, r.Options.Url, body); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return
}

// request method "get" query string
func (r *Request) buildQueryStringRequest() (req *http.Request, err error) {
	var body io.Reader
	data := encodeBodyValues(r.Options.Body)
	var buf bytes.Buffer
	buf.WriteString(data.Encode())
	body = strings.NewReader(buf.String())
	if req, err = http.NewRequest(r.Options.Method, r.Options.Url, body); err != nil {
		return
	}
	return
}

func encodeBodyValues(values map[string]string) url.Values {
	data := url.Values{}
	for k, v := range values {
		data.Add(k, v)
	}
	return data
}

// set go request mime header
func (r *Request) setMIMEHeader(req *http.Request) {
	for k, v := range r.Options.Headers {
		req.Header.Set(k, v)
	}
}

func (r *Request) newGoClient(req *http.Request) (client *http.Client, err error) {
	client = new(http.Client)
	// client cookie
	if err = r.setGoClientCookie(client, req); err != nil {
		return
	}
	// client proxy
	r.setGoClientProxy(client)

	// client timeout
	r.setClientTimeout(client)
	return
}

// client cookie
func (r *Request) setGoClientProxy(client *http.Client) {
	// set client proxy
	if r.Options.ProxyAddr != "" {
		t := new(http.Transport)
		t.Proxy = func(request *http.Request) (*url.URL, error) {
			return url.Parse(r.Options.ProxyAddr)
		}
		client.Transport = t
	}
}

// client proxy
// the proxy eg: "http://ip:prot"
func (r *Request) setGoClientCookie(client *http.Client, req *http.Request) error {
	if jar, err := cookiejar.New(nil); err != nil {
		return err
	} else {
		jar.SetCookies(req.URL, r.Options.Cookies)
		client.Jar = jar
		return nil
	}
}

// client timeout
func (r *Request) setClientTimeout(client *http.Client) {
	client.Timeout = time.Second * time.Duration(conf.GlobalSharedConfig.Http.Timeout)
}

func (r *Request) Do() (*Response, error) {
	var (
		err       error
		client    *http.Client
		goRequest *http.Request
	)
	goRequest, err = r.newGoRequest()
	if err != nil {
		return nil, err
	}
	client, err = r.newGoClient(goRequest)
	if err != nil {
		return nil, err
	}
	if goResponse, err := client.Do(goRequest); err != nil {
		return nil, err
	} else {
		return WrapResponse(r, goResponse), nil
	}
}

func (r *Request) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Request) Unmarshal(bs []byte) (error) {
	return json.Unmarshal(bs, r)
}
