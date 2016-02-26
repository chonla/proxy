package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

type Stub struct {
	List map[string]Recoder `json:"stub"`
}

func recordRequest(req *http.Request) Recoder {
	iBody, err := httputil.DumpRequest(req, true)
	fatal(err)

	return Recoder{
		req: req,
		Request: Inbound{
			Host:   req.Host,
			Path:   req.URL.Path,
			Method: req.Method,
			// Body:   iBody,
			Body: string(iBody),
		},
	}

}

func recordResponse(row Recoder, resp *http.Response) {
	oBody, err := httputil.DumpResponse(resp, true)
	fatal(err)
	row.resp = resp
	row.Response.Status = resp.Status
	row.Response.StatusCode = resp.StatusCode
	// row.Response.Body = oBody
	row.Response.Body = string(oBody)
	row.Name = row.req.Method + "|" + row.req.RequestURI
	data.List[row.Name] = row
}

func writeStub() {
	if arg.Mode == "Record" {
		println()
		println("write stub...")
		b, err := json.Marshal(data)
		fatal(err)

		err = ioutil.WriteFile(arg.StubFileName, b, 0666)
		fatal(err)
	}
}

func (s Stub) ReadFromStub() {
	s.List = make(map[string]Recoder)
	// data.List = make(map[string]Recoder)
	b, err := ioutil.ReadFile(arg.StubFileName)
	if err != nil {
		println("missing stub file", arg.StubFileName)
		return
	}
	// fatal(err)

	err = json.Unmarshal(b, &s)
	// err = json.Unmarshal(b, &data)
	fatal(err)
}

func getResponseFromStub(req *http.Request) *http.Response {
	name := req.Method + "|" + req.RequestURI

	fmt.Printf("name=%s\n", name)
	if row, found := data.List[name]; found {
		// b := row.Response.Body
		b := []byte(row.Response.Body)
		reader := bufio.NewReader(bytes.NewReader(b))
		r, err := http.ReadResponse(reader, req)
		fatal(err)
		return r
	}
	return nil
}
