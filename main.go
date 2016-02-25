package main

// http://stackoverflow.com/questions/31535569/golang-how-to-read-response-body-of-reverseproxy
// curl -verbose -X POST -d @ReadNewCardWS.txt http://athena13:9582/ReadNewCardWS

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kr/pretty"
)

type Arg struct {
	Endpoint     string
	Mode         string
	ProxyPort    string
	StubFileName string
}

var arg Arg

var _ http.RoundTripper = &transport{}

type transport struct {
	http.RoundTripper
}

var data Stub

type Stub struct {
	List map[string]Recoder `json:"stub"`
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
	// Body   []byte
	Body string
}

type Outbound struct {
	Status     string
	StatusCode int
	// Body   []byte
	Body string
}

func (r Recoder) getFromCache(t *transport) (*http.Response, error) {
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
	fmt.Printf("getFromCache: err=%# v, resp=%# v\n\n\n", err, pretty.Formatter(resp))
	return resp, err

}

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

func unproxyURL(req *http.Request) {
	// fmt.Printf("req.RequestURI=%v\nreq.", req.RequestURI, req.Host, req.)
	strURL := req.RequestURI

	// req.Host = req.URL.
	// TODO: add logic convert http to https
	// if req.RequestURI == "http://gliese1dtac-blltxsb.tac.co.th:7844/QueryCDR/CustIntrMgmt/BillEnquiry/SVCBEQryUsageSumm/v2_0/SVCBEQryUsageSumm" {
	// 	strURL = "https://gliese1dtac-blltxsb.tac.co.th:7844/QueryCDR/CustIntrMgmt/BillEnquiry/SVCBEQryUsageSumm/v2_0/SVCBEQryUsageSumm"
	// }
	target, err := url.Parse(strURL)
	fatal(err)
	req.URL = target
}

func recordRequest(req *http.Request) Recoder {
	iBody, err := httputil.DumpRequest(req, true)
	fatal(err)

	return Recoder{
		req: req,
		Request: Inbound{
			Host:   req.Host,
			Path:   req.URL.Path,
			Method: req.Method,
			// Body:   iBody,
			Body: string(iBody),
		},
	}

}

func recordResponse(row Recoder, resp *http.Response) {
	oBody, err := httputil.DumpResponse(resp, true)
	fatal(err)
	row.resp = resp
	row.Response.Status = resp.Status
	row.Response.StatusCode = resp.StatusCode
	// row.Response.Body = oBody
	row.Response.Body = string(oBody)
	row.Name = row.req.Method + "|" + row.req.RequestURI
	data.List[row.Name] = row
}

func fatal(err error) {
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatal(err)
		os.Exit(1)
	}
}

func captureExitProgram() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		writeStub()
		println()
		println("end proxy...")
		os.Exit(1)
	}()
}

func writeStub() {
	if arg.Mode == "Record" {
		println()
		println("write stub...")
		b, err := json.Marshal(data)
		fatal(err)

		err = ioutil.WriteFile(arg.StubFileName, b, 0666)
		fatal(err)
	}
}

func parseArg() {
	flag.StringVar(&arg.Endpoint, "target", "https://curl.haxx.se/", "target URL for reverse proxy")
	flag.StringVar(&arg.ProxyPort, "port", "9000", "proxy running on port. EX: 9000")
	flag.StringVar(&arg.Mode, "mode", "Record", "proxy running mode [Record/Replay]")
	flag.StringVar(&arg.StubFileName, "stubFileName", "stub.txt", "record to file name EX stub.txt")
	flag.Parse()

	println("===========================================")
	println(arg.Mode, "Mode")
	println("===========================================")
	fmt.Printf("arg = %#v\n", arg)
}

func readFromStub() {
	b, err := ioutil.ReadFile(arg.StubFileName)
	fatal(err)

	err = json.Unmarshal(b, &data)
	fatal(err)
}

func getResponseFromStub(req *http.Request) *http.Response {
	name := req.Method + "|" + req.RequestURI

	fmt.Printf("name=%s\n", name)
	if row, found := data.List[name]; found {
		// b := row.Response.Body
		b := []byte(row.Response.Body)
		reader := bufio.NewReader(bytes.NewReader(b))
		r, err := http.ReadResponse(reader, req)
		fatal(err)
		return r
	}
	return nil
}

func getValueByKey(key string, data string) string {
	list := strings.FieldsFunc(data, func(r rune) bool {
		return r == '<' || r == '>'
	})
	for i, s := range list {
		if s == key {
			return list[i+1]
		}
	}
	return ""
}

func main() {
	captureExitProgram()
	println("starting proxy...")
	parseArg()
	data.List = make(map[string]Recoder)

	if arg.Mode == "Replay" {
		readFromStub()
	}

	target, err := url.Parse(arg.Endpoint)
	fatal(err)

	//assign proxy
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &transport{&http.Transport{
		Dial: (&net.Dialer{}).Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS10,
			MaxVersion:         tls.VersionTLS10,
		},
		DisableCompression: false,
		DisableKeepAlives:  true,
	}}

	go serveServerHTTPS(proxy)

	//start proxy
	println("run proxy http")
	log.Fatal(http.ListenAndServe("localhost:"+arg.ProxyPort, nil))
}

type ServerHTTPS struct{}

func (s ServerHTTPS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("https req=%# v\n\n\n", pretty.Formatter(r))
}

func logHttp(w http.ResponseWriter, r *http.Request) {
	println("logHttp")
	//fmt.Printf("log http req=%# v\n\n\n", pretty.Formatter(r))
}

type mx struct{}

var mux map[string]func(http.ResponseWriter, *http.Request)

func serveServerHTTPS(proxy *httputil.ReverseProxy) {
	server := mx{}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = logHttp

	s := &http.Server{
		Addr:    ":9001",
		Handler: server,
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
