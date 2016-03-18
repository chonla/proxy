package main

/*
alt proxy mode
note right of API: host in include list

API-->Proxy: call http to service
Proxy --> *Stub: look in stub
Stub --> API: found in stub
note right of Stub: not found in stub
Stub-->Service: forward request
Service --> Stub: result
Stub --> Proxy: result from stub
Proxy --> API: return result

else direct call
API->Service: call http to service
Service -> API: return result
*/

func main() {
	captureExitProgram()
	println("starting proxy...")

	ReadConfig()
	parseArg()
	ReadFromStub()
	startProxy()
}
