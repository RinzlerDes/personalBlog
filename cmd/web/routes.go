package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// func (app *Application) routes() *http.ServeMux {
func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(app.flags.fileServerAddr))

	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("GET /posts/view/{id}", app.viewHandler)
	mux.HandleFunc("GET /posts/search", app.searchHandler)
	mux.HandleFunc("POST /posts/search", app.searchHandlerProcessForm)
	mux.HandleFunc("GET /posts/insert", app.insertHandler)
	mux.HandleFunc("POST /posts/insert", app.insertHandlerPost)
	// mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// Does the same thing
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		app.fileServerHandler(w, r, fileServer)
	})

	chain := alice.New(secureHeaders, logRequest, app.recoverPanic)
	chain = app.appendSessionHandlers(chain)

	// return mux
	return chain.Then(mux)
	// return app.recoverPanic(logRequest(secureHeaders(mux)))
}

func (app *Application) appendSessionHandlers(chain alice.Chain) alice.Chain {
	// chain = chain.Append(midOne)
	chain = chain.Append(app.sessionManager.LoadAndSave)
	return chain
}
