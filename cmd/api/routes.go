package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func parsePort() (string, error) {

	if len(os.Args) != 2 {
		return "", errors.New("give a port to listen to")
	}
	port := os.Args[1]

	fmt.Printf("Listening on port %s\n", port)
	return fmt.Sprintf(":%s", port), nil
}

func routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ptlist", GetAllTimestamps)

	return mux
}
