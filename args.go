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
}

func parseArg() {
	flag.StringVar(&arg.Endpoint, "target", "https://curl.haxx.se/", "target URL for reverse proxy")
	flag.StringVar(&arg.ProxyPort, "port", "9000", "proxy running on port. EX: 9000")
	flag.StringVar(&arg.Mode, "mode", "Record", "proxy running mode [Record/Replay]")
	flag.StringVar(&arg.StubFileName, "stubFileName", "stub.txt", "record to file name EX stub.txt")
	flag.Parse()

	println("===========================================")
	println(arg.Mode, "Mode")
	println("===========================================")
	fmt.Printf("arg = %#v\n", arg)
}
