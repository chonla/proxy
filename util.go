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

	if inHostList(strings.Join(arg.HttpsList, ","), req.Host) {
		strURL = changeHostToHttps(req.RequestURI)
		fmt.Printf("changeHostToHttps=%v\n", strURL)
	}

	target, err := url.Parse(strURL)
	fatal(err)
	req.URL = target
}

func inHostList(hostList, hostname string) bool {
	index := strings.Index(hostList, hostname)
	return (index >= 0)
}

func changeHostToHttps(endpoint string) string {
	return strings.Replace(endpoint, "http://", "https://", -1)
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
