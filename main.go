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

type Arg struct {
	Endpoint     string
	Mode         string
	ProxyPort    string
	StubFileName string
}

var arg Arg
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
	println("===========================================")
	println("Record Mode")
	println("===========================================")

	var cache *http.Response

	iBody, err := httputil.DumpRequest(req, true)
	// fmt.Printf("req=%#v\n", string(iBody))
	fatal(err)

	row := Recoder{
		req: req,
		Request: Inbound{
			Host:   req.Host,
			Path:   req.URL.Path,
			Method: req.Method,
			Body:   string(iBody),
		},
	}

	if arg.Mode == "Replay" {
		cache = getResponseFromStub()
	}

	fmt.Printf("cache=%v\n\n", cache)
	if cache != nil {
		resp = cache
	} else {
		resp, err = t.RoundTripper.RoundTrip(req)
		fatal(err)
		if err != nil {
			return nil, err
		}
	}
	fmt.Printf("resp=%#v\n\n", resp)

	if arg.Mode == "Record" {
		oBody, err := httputil.DumpResponse(resp, true)
		fatal(err)
		row.resp = resp
		row.Response.Status = resp.Status
		row.Response.StatusCode = resp.StatusCode
		row.Response.Body = string(oBody)
		row.Name = row.req.Method + "|" + row.req.RequestURI + "|" + row.resp.Status + "|"
		data = append(data, row)
		writeStub()
	}

	// // change server to schmerver
	// b, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }
	// err = resp.Body.Close()
	// if err != nil {
	// 	return nil, err
	// }
	// b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	// body := ioutil.NopCloser(bytes.NewReader(b))
	// resp.Body = body
	// resp.ContentLength = int64(len(b))
	// resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	changeServer(resp)
	return resp, nil
}

func main() {
	captureExitProgram()
	println("starting proxy...")
	parseArg()

	target, err := url.Parse(arg.Endpoint)
	fatal(err)

	//assign proxy
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &transport{http.DefaultTransport}

	fmt.Printf("arg = %#v\n", arg)

	//start proxy
	http.Handle("/", proxy)
	// http.HandleFunc("/", report)
	log.Fatal(http.ListenAndServe("localhost:"+arg.ProxyPort, nil))

}

func fatal(err error) {
	if err != nil {
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
}

func parseArg() {
	flag.StringVar(&arg.Endpoint, "target", "http://athena13.tac.co.th:9582/", "target URL for reverse proxy")
	flag.StringVar(&arg.ProxyPort, "port", "9000", "proxy running on port. EX: 9000")
	flag.StringVar(&arg.Mode, "mode", "Record", "proxy running mode [Record/Replay]")
	flag.StringVar(&arg.StubFileName, "stubFileName", "stub.txt", "record to file name EX stub.txt")
	flag.Parse()
}

func getResponseFromStub() *http.Response {
	return nil
}

func changeServer(resp *http.Response) {
	// change server to schmerver
	b, err := ioutil.ReadAll(resp.Body)
	fatal(err)
	// if err != nil {
	// 	return nil, err
	// }
	err = resp.Body.Close()
	fatal(err)
	// if err != nil {
	// 	return nil, err
	// }
	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
}
