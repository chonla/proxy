package main

var arg Arg
var data Stub

func main() {
	captureExitProgram()
	println("starting proxy...")
	parseArg()
	ReadFromStub()
	startProxy()
}
