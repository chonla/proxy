package main

import (
	"flag"
	"fmt"
)

type Arg struct {
	Endpoint     string
	Mode         string
	ProxyPort    string
	StubFileName string
	HttpsList    string

	ConditionList map[string]Condition
}

type Condition struct {
	Url      string
	KeyField string
}

func loadArg() {
	parseArg()
}

func parseArg() {
	flag.StringVar(&arg.Endpoint, "target", "http://1.2.3.4/", "target URL for reverse proxy")
	flag.StringVar(&arg.ProxyPort, "port", "9000", "proxy running on port. EX: 9000")
	flag.StringVar(&arg.Mode, "mode", "Record", "proxy running mode [Record/Replay], default is Record")
	flag.StringVar(&arg.StubFileName, "stubFileName", "stub.json", "record to file name EX stub.json")
	flag.StringVar(&arg.HttpsList, "httpsList", "", "list of https host EX google.com,yahoo.com")
	flag.Parse()

	println("===========================================")
	println(arg.Mode, "Mode")
	println("===========================================")
	fmt.Printf("arg = %#v\n", arg)
}
