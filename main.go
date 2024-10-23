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
	mux.HandleFunc("/page/view", viewHandler)
	mux.HandleFunc("/page/create", createHandler)

	addrString := fmt.Sprintf("localhost:%d", portNumber)

	log.Printf("Starting server on :%v", portNumber)
	log.Fatal(http.ListenAndServe(addrString, mux))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("In home"))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Viewing post"))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not alloweddddd"))
		return
	}
	w.Write([]byte("Creating post"))
}
