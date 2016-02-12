package main

// http://stackoverflow.com/questions/31535569/golang-how-to-read-response-body-of-reverseproxy
// curl -verbose -X POST -d @ReadNewCardWS.txt http://athena13:9582/ReadNewCardWS

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	// "github.com/kr/pretty"
)

var endpoint_url *string
var _ http.RoundTripper = &transport{}
var data []Recoder

type transport struct {
	http.RoundTripper
}

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
	Body   string
}

type Outbound struct {
	Status     string
	StatusCode int
	Body       string
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	// report(resp, r)
	println("===========================================")
	println("Record Mode")
	println("===========================================")

	iBody, err := httputil.DumpRequest(req, true)

	resp, err = t.RoundTripper.RoundTrip(req)

	oBody, err := httputil.DumpResponse(resp, true)

	row := Recoder{
		req:  req,
		resp: resp,
		Request: Inbound{
			Host:   req.Host,
			Path:   req.URL.Path,
			Method: req.Method,
			Body:   string(iBody),
		},
		Response: Outbound{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Body:       string(oBody),
		},
	}
	row.Name = row.req.Method + "|" + row.req.RequestURI + "|" + row.resp.Status + "|"
	data = append(data, row)

	writeStub()

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
	captureEnd()
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

func captureEnd() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		println()
		println("end proxy...")

		os.Exit(1)
	}()

}

func writeStub() {
	println("write stub...")
	b, err := json.Marshal(data)
	fatal(err)

	fmt.Printf("data=%s", string(b))
	err = ioutil.WriteFile("stub.txt", b, 0666)
	fatal(err)

	// fmt.Printf("data=%v\n", pretty.Formatter(data))

}
