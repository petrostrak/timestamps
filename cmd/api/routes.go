package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func parsePort() string {

	if len(os.Args) != 2 {
		fmt.Println("give a port to listen to")
		os.Exit(0)
	}
	port := os.Args[1]

	fmt.Printf("Listening on port %s\n", port)
	return ":" + port
}

func routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ptlist", GetAllTimestamps)

	return mux
}

func startApp() {
	err := http.ListenAndServe(parsePort(), routes())
	if err != nil {
		log.Fatal(err)
	}
}
