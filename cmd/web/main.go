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
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/posts/view", viewHandler)
	mux.HandleFunc("/posts/create", createHandler)
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		fileServerHandler(w, r, fileServer)
	})

	addrString := fmt.Sprintf("localhost:%d", portNumber)

	log.Printf("Starting server on :%v", portNumber)
	log.Fatal(http.ListenAndServe(addrString, mux))
}
