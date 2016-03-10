package main

func main() {
	captureExitProgram()
	println("starting proxy...")

	ReadConfig()
	parseArg()
	ReadFromStub()
	startProxy()
}
