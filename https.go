package main

// http://stackoverflow.com/questions/31535569/golang-how-to-read-response-body-of-reverseproxy
// curl -verbose -X POST -d @ReadNewCardWS.txt http://athena13:9582/ReadNewCardWS

import (
	"crypto/tls"
	// "fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var _ http.RoundTripper = &Transport{}

type Transport struct {
	http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	row := newRecoder(req)
	unproxyURL(req)
	resp, err = row.getFromCache(t)
	return
}

func startProxy() {
	target, err := url.Parse(arg.Endpoint)
	fatal(err)

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &Transport{&http.Transport{
		Dial: (&net.Dialer{}).Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS10,
			MaxVersion:         tls.VersionTLS10,
		},
		DisableCompression: false,
		DisableKeepAlives:  true,
	}}

	//start proxy
	println("runing proxy http...")
	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe("localhost:"+arg.ProxyPort, nil))
}
