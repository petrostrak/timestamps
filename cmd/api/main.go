package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port, err := parsePort()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	handler := routes()

	err = http.ListenAndServe(port, handler)
	if err != nil {
		log.Fatal(err)
	}
}
