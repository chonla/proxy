package main

var arg Arg
var data Stub

func main() {
	captureExitProgram()
	println("starting proxy...")
	data.List = make(map[string]Recoder)
	parseArg()
	readFromStub()
	startProxy()
}
