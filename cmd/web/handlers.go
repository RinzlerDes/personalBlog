package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err)
		http.Error(w, "internal errorrrr", http.StatusInternalServerError)
	}

	err = t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal errorrrr", http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		str := fmt.Sprintf("Post %v does not exist\n", id)
		http.Error(w, str, http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Viewing post %v\n", id)
	//w.Write([]byte("Viewing post\n"))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		// System headers if need to be suppressed, need to be done manually
		//w.Header()["Date"] = nil
		// Does the job of the two commented out lines below
		http.Error(w, "Method not allowedddd", http.StatusMethodNotAllowed)
		// w.WriteHeader(http.StatusMethodNotAllowed)
		// w.Write([]byte("Method not alloweddddd\n"))
		return
	}
	w.Write([]byte("Creating post"))
}
