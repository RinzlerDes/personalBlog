package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type CommandLineFlags struct {
	addr           string
	fileServerAddr string
}

func main() {
	logErr := log.New(os.Stderr, "ERRORR:", log.Lshortfile|log.Ltime|log.Ldate)
	logInfo := log.New(os.Stdout, "INFOO:", log.Lshortfile|log.Ltime|log.Ldate)

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

	logInfo.Printf("Starting server on %s", commandLineFlags.addr)
	logErr.Fatal(http.ListenAndServe(commandLineFlags.addr, mux))
}

func (flags *CommandLineFlags) getCommandLineFlags() {
	flag.StringVar(&flags.addr, "addr", "localhost:8080", "HTTP network address")
	flag.StringVar(&flags.fileServerAddr, "fileServerAddr", "./ui/static", "Path to static assets")
	flag.Parse()
}
