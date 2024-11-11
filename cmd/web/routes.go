package main

import "net/http"

func (app *Application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(app.flags.fileServerAddr))

	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/posts/view", app.viewHandler)
	mux.HandleFunc("/posts/create", app.createHandler)
	mux.HandleFunc("/posts/search", app.searchHandler)
	mux.HandleFunc("/posts/insert", app.insertHandler)
	// mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// Does the same thing
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		app.fileServerHandler(w, r, fileServer)
	})

	return mux
}
