package main

import (
	"os"
	"testing"
)

func TestParsePorts(t *testing.T) {

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"no port", "", ":"},
		{"with port", "8080", ":8080"},
	}

	for _, e := range testCases {
		os.Args = []string{"timestamps", e.input}
		res := parsePort()
		if e.expected != res {
			t.Errorf("Got %s but wanted %s", res, e.expected)
		}
	}
}
