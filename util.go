package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func unproxyURL(req *http.Request) {
	fmt.Printf("req.RequestURI=%v\n", req.RequestURI)
	strURL := req.RequestURI

	// TODO: add logic convert http to https
	if req.RequestURI == "http://gliese1dtac-blltxsb.tac.co.th:7844/QueryCDR/CustIntrMgmt/BillEnquiry/SVCBEQryUsageSumm/v2_0/SVCBEQryUsageSumm" {
		strURL = "https://gliese1dtac-blltxsb.tac.co.th:7844/QueryCDR/CustIntrMgmt/BillEnquiry/SVCBEQryUsageSumm/v2_0/SVCBEQryUsageSumm"
	}

	if inHostList(arg.HttpsList, req.Host) {

	}

	target, err := url.Parse(strURL)
	fatal(err)
	req.URL = target
}

func inHostList(hostList, hostname string) bool {
	index := strings.Index(hostList, hostname)
	return (index >= 0)
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
		data.WriteStub()

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
