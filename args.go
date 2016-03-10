package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/kr/pretty"
)

type Url string
type KeyField string

type Arg struct {
	Endpoint     string
	Mode         string
	ProxyPort    string
	StubFileName string
	HttpsList    string

	IncludeList map[Url]KeyField
}

var arg Arg

func loadArg() {
	parseArg()
}

func parseArg() {
	ReadConfig()

	flag.StringVar(&arg.ProxyPort, "port", "9000", "proxy running on port. EX: 9000")
	flag.StringVar(&arg.Mode, "mode", "Record", "proxy running mode [Record/Replay], default is Record")
	flag.Parse()

	println("===========================================")
	println(arg.Mode, "Mode")
	println("===========================================")
	fmt.Printf("config %# v \n\n", pretty.Formatter(arg))
}

func ReadConfig() {
	filename := "proxy.json"
	arg.Endpoint = "http://1.2.3.4/"
	println("loading config file", filename)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		println("missing config file", filename)
		return
	}
	err = json.Unmarshal(b, &arg)
	fatal(err)
}
