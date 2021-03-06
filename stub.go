package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Stub struct {
	List map[string]Recoder `json:"stub"`
}

var data Stub

func (s Stub) WriteStub() {
	if arg.Mode == "Record" {
		println()
		println("write stub file ...", arg.StubFileName)
		b, err := json.Marshal(s)
		fatal(err)

		err = ioutil.WriteFile(arg.StubFileName, b, 0666)
		fatal(err)
	}
}

func ReadFromStub() {
	data.List = make(map[string]Recoder)
	b, err := ioutil.ReadFile(arg.StubFileName)
	if err != nil {
		println("missing stub file", arg.StubFileName)
		return
	}

	err = json.Unmarshal(b, &data)
	fmt.Printf("CACHE: loaded current cache %v record\n\n", len(data.List))
}

func (s Stub) FindInCache(r Recoder) *http.Response {
	name := generateKey(r.Request)
	if row, found := s.List[name]; found {
		b := []byte(row.Response.Body)
		reader := bufio.NewReader(bytes.NewReader(b))
		resp, err := http.ReadResponse(reader, r.req)
		fatal(err)
		return resp
	}
	return nil
}
