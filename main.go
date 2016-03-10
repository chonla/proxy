package main

func main() {
	captureExitProgram()
	println("starting proxy...")
	parseArg()
	ReadFromStub()
	startProxy()
}
