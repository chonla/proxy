package main

import "testing"

func TestFoundInHostList(t *testing.T) {
	var testcases = []struct {
		name     string
		expected bool
		list     string
		search   string
	}{
		{"match", true, "google.com,yahoo.com", "google.com"},
		{"match", true, "google.com,yahoo.com", "yahoo.com"},
		{"not match", false, "google.com,yahoo.com", "bing.com"},
		{"not match", false, "", "bing.com"},
	}

	for _, testcase := range testcases {
		if inHostList(testcase.list, testcase.search) != testcase.expected {
			t.Error("fail case ", testcase.name)
		}
	}
}
