package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	//app.logErr.Output(2, trace)
	logErr.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (flags *CommandLineFlags) getCommandLineFlags() {
	flag.StringVar(&flags.addr, "addr", "localhost:8080", "HTTP network address")
	flag.StringVar(&flags.fileServerAddr, "fileServerAddr", "./ui/static", "Path to static assets")
	flag.Parse()
}
