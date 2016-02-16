package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		fmt.Printf("HEADER: %v=%v\n", n, v)
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
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
