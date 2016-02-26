package main

// http://stackoverflow.com/questions/31535569/golang-how-to-read-response-body-of-reverseproxy
// curl -verbose -X POST -d @ReadNewCardWS.txt http://athena13:9582/ReadNewCardWS

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kr/pretty"
)

var _ http.RoundTripper = &Transport{}

type Transport struct {
	http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	fmt.Printf("req=%# v\n\n\n", pretty.Formatter(req))
	row := recordRequest(req)
	unproxyURL(req)
	resp, err = row.getFromCache(t)
	if err != nil {
		recordResponse(row, resp)
	}

	return resp, nil
}

func startProxy() {
	//assign proxy
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

	// go ServeTCP()

	//start proxy
	println("run proxy http")
	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe("localhost:"+arg.ProxyPort, nil))

}

// func servHTTPS() {
// 	//start https proxy
// 	println("run proxy https")
// 	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
// 	if err != nil {
// 		log.Fatalf("server: loadkeys: %s", err)
// 	}
// 	config := tls.Config{
// 		Certificates:       []tls.Certificate{cert},
// 		InsecureSkipVerify: true,
// 		MinVersion:         tls.VersionTLS10,
// 		MaxVersion:         tls.VersionTLS10,
// 	}
// 	config.Rand = rand.Reader
// 	service := ":8000"
// 	listener, err := tls.Listen("tcp", service, &config)
// 	if err != nil {
// 		log.Fatalf("server: listen: %s", err)
// 	}
// 	fatal(http.Serve(listener, nil))
// 	// fatal(http.ListenAndServeTLS("localhost:8000", "cert.pem", "key.pem", nil))

// }

// func ServeTCP() {
// 	println("run proxy tcp 8000")
// 	l, err := net.Listen("tcp", ":8000")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer l.Close()
// 	for {
// 		// Wait for a connection.
// 		conn, err := l.Accept()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// Handle the connection in a new goroutine.
// 		// The loop then returns to accepting, so that
// 		// multiple connections may be served concurrently.
// 		go func(c net.Conn) {
// 			// Echo all incoming data.

// 			fmt.Printf("c=%# v", c)

// 			println("hello")
// 			io.Copy(c, c)

// 			// Shut down the connection.
// 			c.Close()
// 		}(conn)
// 	}
// }
