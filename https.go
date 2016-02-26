package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	// "net/http/httputil"
	"crypto/rand"
	"io"
	"net"

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

func servHTTPS() {
	//start https proxy
	println("run proxy https")
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS10,
	}
	config.Rand = rand.Reader
	service := ":8000"
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	fatal(http.Serve(listener, nil))
	// fatal(http.ListenAndServeTLS("localhost:8000", "cert.pem", "key.pem", nil))

}

func ServeTCP() {
	println("run proxy tcp 8000")
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.

			fmt.Printf("c=%# v", c)

			println("hello")
			io.Copy(c, c)

			// Shut down the connection.
			c.Close()
		}(conn)
	}
}

/*package main

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
*/
