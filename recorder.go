package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	// "github.com/kr/pretty"
)

type Recoder struct {
	Name     string `json:"-"`
	Request  Inbound
	Response Outbound
	req      *http.Request  `json:"-"`
	resp     *http.Response `json:"-"`
}

type Inbound struct {
	URI      string `json:"-"`
	Host     string
	Path     string
	Method   string
	Body     []byte
	BodyText string `json:"-"`
}

type Outbound struct {
	Status     string
	StatusCode int
	Body       []byte
	BodyText   string `json:"-"`
}

func recordRequest(req *http.Request) Recoder {
	iBody, err := httputil.DumpRequest(req, true)
	fatal(err)

	fmt.Printf("\n\nPOST BODY: %v \n\n", string(iBody))

	return Recoder{
		req: req,
		Request: Inbound{
			Host:     req.Host,
			Path:     req.URL.Path,
			Method:   req.Method,
			Body:     iBody,
			BodyText: string(iBody),
		},
	}
}

func recordResponse(row Recoder, resp *http.Response) {
	oBody, err := httputil.DumpResponse(resp, true)
	fatal(err)
	row.resp = resp
	row.Response.Status = resp.Status
	row.Response.StatusCode = resp.StatusCode
	row.Response.Body = oBody
	row.Response.BodyText = string(oBody)
	row.Name = row.req.Method + "|" + row.req.RequestURI

	data.List[row.Name] = row
	fmt.Printf("CACHE: added current cache %v record\n", len(data.List))
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
