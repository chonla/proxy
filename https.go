package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/kr/pretty"
)

type transport struct {
	http.RoundTripper
}
type ServerHTTPS struct{}
type mx struct{}

var mux map[string]func(http.ResponseWriter, *http.Request)

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	fmt.Printf("RoundTrip req=%# v\n\n\n", pretty.Formatter(req))
	row := recordRequest(req)
	unproxyURL(req)
	fmt.Printf("unproxyURL req=%# v\n\n\n", pretty.Formatter(req))

	resp, err = row.getFromCache(t)
	if err == nil {
		recordResponse(row, resp)
	}

	return resp, nil
}

func (s ServerHTTPS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("https req=%# v\n\n\n", pretty.Formatter(r))
}

func logHttp(w http.ResponseWriter, r *http.Request) {
	println("logHttp")
	fmt.Printf("log http req=%# v\n\n\n", pretty.Formatter(r))
}

func serveServerHTTPS(proxy *httputil.ReverseProxy) {
	server := mx{}
	// mux = make(map[string]mx)
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["//www.google.com:443"] = logHttp

	s := &http.Server{
		Addr:    ":9001",
		Handler: &server,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS10,
			MaxVersion:         tls.VersionTLS10,
		},
	}
	log.Fatal(s.ListenAndServe())
}

func (m mx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		println("match")
		h(w, r)
		return
	}
	fmt.Printf("My server: %v\n", r.URL.String())
	fmt.Printf("mux.ServeHTTP req=%# v\n\n\n", pretty.Formatter(r))
}
