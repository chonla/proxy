package main

import (
	"fmt"
	"net/http"

	"github.com/kr/pretty"
)

type Recoder struct {
	Name     string
	Request  Inbound
	Response Outbound
	req      *http.Request  `json:"-"`
	resp     *http.Response `json:"-"`
}

type Inbound struct {
	URI    string
	Host   string
	Path   string
	Method string
	// Body   []byte
	Body string
}

type Outbound struct {
	Status     string
	StatusCode int
	// Body   []byte
	Body string
}

func (r Recoder) getFromCache(t *Transport) (*http.Response, error) {
	var cache *http.Response
	if arg.Mode == "Replay" {
		cache = getResponseFromStub(r.req)
	}

	if cache != nil {
		fmt.Printf("cache hit=%#v\n", pretty.Formatter(cache))
		return cache, nil
	}
	println("cache miss, call http")
	resp, err := t.RoundTripper.RoundTrip(r.req)
	fmt.Printf("err=%# v, resp=%# v\n\n\n", err, pretty.Formatter(resp))
	return resp, err

}
