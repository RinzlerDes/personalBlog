package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	portNumber uint = 8080
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/posts/view", viewHandler)
	mux.HandleFunc("/posts/create", createHandler)

	addrString := fmt.Sprintf("localhost:%d", portNumber)

	log.Printf("Starting server on :%v", portNumber)
	log.Fatal(http.ListenAndServe(addrString, mux))
}
