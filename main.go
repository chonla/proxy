package main

var arg Arg
var data Stub

func main() {
	captureExitProgram()
	println("starting proxy...")
	parseArg()
	data.List = make(map[string]Recoder)
	data.ReadFromStub()
	startProxy()
}
