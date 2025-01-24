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
	mux.HandleFunc("POST /posts/search", app.searchHandlerPost)
	mux.HandleFunc("GET /posts/insert", app.insertHandler)
	mux.HandleFunc("POST /posts/insert", app.insertHandlerPost)
	// mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// Does the same thing
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		app.fileServerHandler(w, r, fileServer)
	})

	// chain := alice.New(midThree, secureHeaders, logRequest, app.recoverPanic)
	// chain = app.appendSessionHandlers(chain,
	// 	app.sessionManager.LoadAndSave,
	// )
	chain := alice.Chain{}
	chain = app.appendSessionHandlers(chain,
		app.sessionManager.LoadAndSave,
		secureHeaders,
		logRequest,
		app.recoverPanic,
	)

	// return mux
	return chain.Then(mux)
	// return app.recoverPanic(logRequest(secureHeaders(mux)))
}

// func (app *Application) appendSessionHandlers(handlers ...alice.Constructor) alice.Chain {
func (app *Application) appendSessionHandlers(chain alice.Chain, handlers ...alice.Constructor) alice.Chain {
	// chain := alice.Chain{}
	for _, handler := range handlers {
		chain = chain.Append(handler)
	}
	return chain
}
