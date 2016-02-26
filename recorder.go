package main

import (
	"fmt"
	"net/http"

	// "github.com/kr/pretty"
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
	println()
	cache := data.FindInCache(r.req)
	if cache != nil {
		fmt.Printf("CACHE: hit current cache %v record\n", len(data.List))
		return cache, nil
	}
	fmt.Printf("CACHE: miss current cache %v record, call http\n", len(data.List))

	resp, err := t.RoundTripper.RoundTrip(r.req)
	if err == nil {
		recordResponse(r, resp)
	}

	return resp, err
}
