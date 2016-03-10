package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/kr/pretty"
)

type Condition map[string]string

type Arg struct {
	Endpoint  string `json:"-"`
	Mode      string `json:"-"`
	ProxyPort string `json:"-"`

	StubFileName string
	HttpsList    []string
	IncludeList  Condition
}

var arg Arg

func parseArg() {
	flag.StringVar(&arg.ProxyPort, "port", "9000", "proxy running on port. EX: 9000")
	flag.StringVar(&arg.Mode, "mode", "Record", "proxy running mode [Record/Replay], default is Record")
	flag.Parse()

	// WriteConfig()
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

// func WriteConfig() {
// 	filename := "proxy.json"
// 	println("write config file", filename)
// 	arg.HttpsList = append(arg.HttpsList, "gliese1dtac-blqrysb.tac.co.th:7834")
// 	arg.HttpsList = append(arg.HttpsList, "gliese1-blqrysb.tac.co.th:7832")

// 	b, err := json.Marshal(arg)
// 	fatal(err)

// 	err = ioutil.WriteFile(filename, b, 0666)
// 	fatal(err)

// }
