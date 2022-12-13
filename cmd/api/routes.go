package main

import (
	"fmt"
	"net/http"
	"os"
)

var (
	mux *http.ServeMux
)

func init() {
	mux = http.NewServeMux()
}

func parsePort() string {

	if len(os.Args) != 2 {
		fmt.Println("give a port to listen to")
		os.Exit(0)
	}
	port := os.Args[1]

	fmt.Printf("Listening on port %s\n", port)
	return ":" + port
}

func startApp() {
	mux.HandleFunc("/ptlist", GetAllTimestamps)

	if err := http.ListenAndServe(parsePort(), mux); err != nil {
		panic(err)
	}
}
