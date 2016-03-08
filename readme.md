cb41e82 add func changeHostToHttps
8c6ac72 add flag httpsList
6ac759c Fix case EOF when read from stub file
0f89a1f change logic when write response

2f49a13 change getResponseFromStub to Stub.FindInCache
a642c8c move writeStub to method Stub.WriteStub
e4944c1 move readFromStub to method Stub.ReadFromStub
18b225b remove code comment
49c1655 refactor to func startProxy
132be0c comment unuse func and check stub file after load
4960183 remove code comment
60100e3 FIX:run proxy fail and seperate file again
5216b96 move txt file to sampleFile
33e354b change default mode to replay, record only record mode
dba942f change default stub.txt to stub.json
99977aa seperate function to other file
ae801e3 remove unused file
2a1641d remove servHTTPS, serveTCP func
7daf1d5 add file certificate
93104aa spike listen https via tcp
ff523df add logic when call http error do not record response
8934cdf refactor remove changeServer
1971e9a change test file
c3e0bcd add func getValueByKey
460c531 add proxy_test
4732afb add func unproxyURL
b76fc4a add logic get data from stub
8f85465 change json structure
69dac94 re-arrange code structure
8cad199 add logic Record, Replay
2b1b35d add unused.go and backup old stub
6044dca add new flag target,port,mode,stubFileName
0c3a955 add write to stub file
0d0dbec add capture ctrl+c when exit
821f491 add soap message data file
4fe51f2 add debug round trip
