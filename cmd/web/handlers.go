package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/nav.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		// app.logErr.Println(err)
		// http.Error(w, "internal errorrrr", http.StatusInternalServerError)
		app.serveError(w, err)
		return
	}

	err = t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *Application) viewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		str := fmt.Sprintf("%s\nPost %v does not exist\n", http.StatusText(http.StatusNotFound), id)
		http.Error(w, str, http.StatusNotFound)
		app.logErr.Printf("%v id=%v", err, id)
		return
	}
	fmt.Fprintf(w, "Viewing post %v\n", id)
	//w.Write([]byte("Viewing post\n"))
}

func (app *Application) createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		// System headers if need to be suppressed, need to be done manually
		//w.Header()["Date"] = nil
		// Does the job of the two commented out lines below
		// w.WriteHeader(http.StatusMethodNotAllowed)
		// w.Write([]byte("Method not alloweddddd\n"))
		http.Error(w, "Method not allowedddd", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Creating post"))
}

func (app *Application) fileServerHandler(w http.ResponseWriter, r *http.Request, h http.Handler) {

	orig := r.URL.Path
	strippedPath := r.URL.Path[len("/static/"):] // Get the path after stripping
	app.logInfo.Printf("Before: %s\nFile path after StripPrefix: %s", orig, strippedPath)

	// Adjust the request URL path to match the file server's expectation
	r.URL.Path = strippedPath // Set the adjusted path for the file server

	// Serve the file using the adjusted request
	h.ServeHTTP(w, r)
}
