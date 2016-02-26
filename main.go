package main

// http://stackoverflow.com/questions/31535569/golang-how-to-read-response-body-of-reverseproxy
// curl -verbose -X POST -d @ReadNewCardWS.txt http://athena13:9582/ReadNewCardWS

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var arg Arg
var data Stub

func main() {
	captureExitProgram()
	println("starting proxy...")
	parseArg()
	data.List = make(map[string]Recoder)

	// if arg.Mode == "Replay" {
	readFromStub()
	// }

	target, err := url.Parse(arg.Endpoint)
	fatal(err)

	//assign proxy
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

	// go ServeTCP()

	//start proxy
	println("run proxy http")
	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe("localhost:"+arg.ProxyPort, nil))
}
