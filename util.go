package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	// "github.com/kr/pretty"
)

func unproxyURL(req *http.Request) {
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
	return searchString(hostList, hostname)
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
		if s == key && len(list) > i && list[i+1] != "/"+key {
			return list[i+1]
		}
	}
	return ""
}

func generateKey(req Inbound) string {
	conditionField := getConditionField(req.Host+req.Path, arg.IncludeList)
	conditionValue := getConditionValue(conditionField, byteToStr(req.Body))
	return req.Method + "|" + req.Host + req.Path + "|" + conditionValue
}

func getConditionField(endpoint string, fieldList Condition) string {
	if list, found := fieldList[endpoint]; found {
		return list
	}
	return ""
}

func getConditionValue(key, data string) string {
	var result []string
	list := strings.Split(key, ",")
	for _, value := range list {
		if v := getValueByKey(value, data); v != "" {
			result = append(result, v)
		}
	}
	return strings.Join(result, "|")
}

func byteToStr(data []byte) string {
	return fmt.Sprintf("%s", data)
}

func isRecordMode() bool {
	return arg.Mode == "Record"
}

func isReplayMode() bool {
	return arg.Mode == "Replay"
}

func foundIncludeList(r Recoder) bool {
	endpoint := r.Request.Host + r.Request.Path
	return exactMatch(endpoint) || wildcardMatch(endpoint)
}

func exactMatch(endpoint string) bool {
	_, foundExact := arg.IncludeList[endpoint]
	return foundExact
}

func wildcardMatch(endpoint string) bool {
	for key, _ := range arg.IncludeList {
		if isRegularExpression(key) {
			match, err := regexp.MatchString(key, endpoint)
			fatal(err)

			if match {
				return true
			}
		}
	}
	return false
}

func isRegularExpression(expr string) bool {
	if !hasAsteriskCharactor(expr) {
		return false
	}

	_, err := regexp.Compile(expr)
	if err == nil {
		return true
	}
	fatal(err)
	return false
}

func hasAsteriskCharactor(expr string) bool {
	return searchString(expr, "*")
}

func searchString(str, search string) bool {
	index := strings.Index(str, search)
	return (index >= 0)
}
