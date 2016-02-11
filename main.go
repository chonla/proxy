package main

// http://stackoverflow.com/questions/31535569/golang-how-to-read-response-body-of-reverseproxy

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"github.com/kr/pretty"
)

var endpoint_url *string
var _ http.RoundTripper = &transport{}

type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	// report(resp, r)
	println("===========================================")
	println("in RoundTrip")

	fmt.Printf("req=%# v\n", pretty.Formatter(req))
	// fmt.Printf("req=%#v\n", req)
	resp, err = t.RoundTripper.RoundTrip(req)

	println("===========================================")
	println("response")
	// fmt.Printf("err=%# v RoundTrip=%#v\n", err, resp)
	fmt.Printf("err=%# v RoundTrip=%# v\n", err, pretty.Formatter(resp))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return resp, nil
}

func main() {
	println("starting proxy...")
	endpoint_url = flag.String("target", "http://athena13:9582/", "target URL for reverse proxy")
	flag.Parse()

	fmt.Sprintf("endpoint_url=%s\n", endpoint_url)

	// ...
	target, err := url.Parse("http://athena13.tac.co.th:9582/")
	fatal(err)

	//assign proxy
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &transport{http.DefaultTransport}

	//start proxy
	http.Handle("/", proxy)
	// http.HandleFunc("/", report)
	log.Fatal(http.ListenAndServe("localhost:9000", nil))

}

func report(w http.ResponseWriter, r *http.Request) {
	println("Starting request ....")
	uri := *endpoint_url + r.RequestURI

	fmt.Println(r.Method + ": " + uri)

	rr, err := http.NewRequest(r.Method, uri, r.Body)
	fatal(err)
	copyHeader(r.Header, &rr.Header)

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		fatal(err)
		fmt.Printf("Body: %v\n", string(body))
	}

	// rr, err := http.NewRequest(r.Method, uri, r.Body)
	// fatal(err)
	// copyHeader(r.Header, &rr.Header)

	// Create a client and query the target
	var transport http.Transport
	resp, err := transport.RoundTrip(rr)
	fatal(err)

	println("===========================================")
	println()

	println("Starting response ....")
	fmt.Printf("Resp-Headers: %v\n", resp.Header)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fatal(err)

	dH := w.Header()
	copyHeader(resp.Header, &dH)
	dH.Add("Requested-Host", rr.Host)

	w.Write(body)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		fmt.Printf("HEADER: %v=%v\n", n, v)
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}
