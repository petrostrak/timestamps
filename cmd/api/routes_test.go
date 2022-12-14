package main

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
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

func TestRoutes(t *testing.T) {

	srv := &http.Server{Addr: ":8080", Handler: routes()}

	go func() {
		time.Sleep(1 * time.Second)
		srv.Shutdown(context.Background())
	}()

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		t.Error("unexpected error:", err)
	}

}
