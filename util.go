package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

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
		if arg.Mode == "Record" {
			writeStub()
		}
		println()
		println("end proxy...")
		os.Exit(1)
	}()
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

func unproxyURL(req *http.Request) {
	println("unproxyURL")
	// fmt.Printf("req.RequestURI=%v\nreq.", req.RequestURI, req.Host, req.)
	strURL := req.RequestURI

	// TODO: add logic convert http to https
	// if req.RequestURI == "http://gliese1dtac-blltxsb.tac.co.th:7844/QueryCDR/CustIntrMgmt/BillEnquiry/SVCBEQryUsageSumm/v2_0/SVCBEQryUsageSumm" {
	// 	strURL = "https://gliese1dtac-blltxsb.tac.co.th:7844/QueryCDR/CustIntrMgmt/BillEnquiry/SVCBEQryUsageSumm/v2_0/SVCBEQryUsageSumm"
	// }
	target, err := url.Parse(strURL)
	fatal(err)
	req.URL = target
}
