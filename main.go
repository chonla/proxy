package main

// http://stackoverflow.com/questions/31535569/golang-how-to-read-response-body-of-reverseproxy
// curl -verbose -X POST -d @ReadNewCardWS.txt http://athena13:9582/ReadNewCardWS

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	// "io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
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
	List map[string]Recoder `json:"list"`
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
	Body   []byte
}

type Outbound struct {
	Status     string
	StatusCode int
	Body       []byte
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	var cache *http.Response

	iBody, err := httputil.DumpRequest(req, true)
	fatal(err)

	row := Recoder{
		req: req,
		Request: Inbound{
			Host:   req.Host,
			Path:   req.URL.Path,
			Method: req.Method,
			Body:   iBody,
		},
	}

	if arg.Mode == "Replay" {
		cache = getResponseFromStub(req)
	}

	if cache != nil {
		resp = cache
		fmt.Printf("cache hit=%v\n", pretty.Formatter(cache))
	} else {
		println("call http")
		resp, err = t.RoundTripper.RoundTrip(req)
		fatal(err)
		if err != nil {
			return nil, err
		}
	}

	if arg.Mode == "Record" {
		oBody, err := httputil.DumpResponse(resp, true)
		fatal(err)
		row.resp = resp
		row.Response.Status = resp.Status
		row.Response.StatusCode = resp.StatusCode
		row.Response.Body = oBody
		row.Name = req.Method + "|" + req.RequestURI
		data.List[row.Name] = row
	}

	changeServer(resp)
	return resp, nil
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
	proxy.Transport = &transport{http.DefaultTransport}

	//start proxy
	http.Handle("/", proxy)
	// http.HandleFunc("/", report)
	log.Fatal(http.ListenAndServe("localhost:"+arg.ProxyPort, nil))
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

		// fmt.Printf("data=%s", string(b))
		err = ioutil.WriteFile("stub.txt", b, 0666)
		fatal(err)
	}
}

func parseArg() {
	flag.StringVar(&arg.Endpoint, "target", "http://athena13.tac.co.th:9582/", "target URL for reverse proxy")
	flag.StringVar(&arg.ProxyPort, "port", "9000", "proxy running on port. EX: 9000")
	flag.StringVar(&arg.Mode, "mode", "Record", "proxy running mode [Record/Replay]")
	flag.StringVar(&arg.StubFileName, "stubFileName", "stub.txt", "record to file name EX stub.txt")
	flag.Parse()

	println("===========================================")
	println(arg.Mode, " Mode")
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

	if row, found := data.List[name]; found {
		b := row.Response.Body
		reader := bufio.NewReader(bytes.NewReader(b))
		r, err := http.ReadResponse(reader, req)
		fatal(err)
		return r
	}
	return nil
}

func changeServer(resp *http.Response) {
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
