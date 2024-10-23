package main

import (
	"flag"
	"log"
	"net/http"
)

type CommandLineFlags struct {
	addr           string
	fileServerAddr string
}

func main() {
	var commandLineFlags CommandLineFlags
	commandLineFlags.getCommandLineFlags()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(commandLineFlags.fileServerAddr))

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/posts/view", viewHandler)
	mux.HandleFunc("/posts/create", createHandler)
	// mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// Does the same thing
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		fileServerHandler(w, r, fileServer)
	})

	log.Printf("Starting server on %s", commandLineFlags.addr)
	log.Fatal(http.ListenAndServe(commandLineFlags.addr, mux))
}

func (flags *CommandLineFlags) getCommandLineFlags() {
	addr := flag.String("addr", "localhost:8080", "HTTP network address")

	flag.Parse()
	flags.addr = *addr
}
