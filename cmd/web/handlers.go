package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"personalBlog/internal/models"
	"strconv"
)

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	posts, err := app.postModel.Latest(5)
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/nav.html",
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = t.ExecuteTemplate(w, "base", posts)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *Application) viewHandler(w http.ResponseWriter, r *http.Request) {
	// Get ID from url
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		str := fmt.Sprintf("%s\nPost %v does not exist\n", http.StatusText(http.StatusNotFound), id)
		http.Error(w, str, http.StatusNotFound)
		logErr.Printf("%v id=%v", err, id)
		return
	}

	// Get the post from DB using the ID
	post, err := app.postModel.Get(uint(id))
	if err != nil {
		// No matching record error
		if errors.Is(err, models.ErrNoRecord) {
			logErr.Println("post not found: ", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Internal server error
		http.Error(w, err.Error(), http.StatusConflict)
		app.serverError(w, err)
		return
	}

	// Render html templates
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/view.html",
		"./ui/html/partials/nav.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = t.ExecuteTemplate(w, "base", post)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the retrieved row on the webpage
	// str := fmt.Sprintf("Retrieved Post\n" + post.String())
	// fmt.Fprintf(w, str)
	//
	// fmt.Fprintf(w, "Viewing post %v\n", id)
	// w.Write([]byte("Viewing post\n"))
}

func (app *Application) createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowedddd", http.StatusMethodNotAllowed)
		return
	}

	title := r.URL.Query().Get("title")
	content := r.URL.Query().Get("content")

	if title == "" || content == "" {
		app.serverError(w, fmt.Errorf("Title nor content can be empty"))
		return
	}

	newPost := models.Post{
		Title:   title,
		Content: content,
	}

	err := app.postModel.Insert(&newPost)
	if err != nil {
		app.serverError(w, fmt.Errorf("could not insert"))
		return
	}

	fmt.Fprintf(w, "Created a new post!!!\nID      %d\nTitle       %s\nContent     %s\nCreated     %v\n",
		newPost.ID, newPost.Title, newPost.Content, newPost.Created)

	w.Write([]byte("Creating post"))
}

func (app *Application) searchHandler(w http.ResponseWriter, r *http.Request) {
	// Render html templates
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/postSearch.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *Application) fileServerHandler(w http.ResponseWriter, r *http.Request, h http.Handler) {
	orig := r.URL.Path
	strippedPath := r.URL.Path[len("/static/"):] // Get the path after stripping
	logInfo.Printf("Before: %s\nFile path after StripPrefix: %s", orig, strippedPath)

	// Adjust the request URL path to match the file server's expectation
	r.URL.Path = strippedPath // Set the adjusted path for the file server

	// Serve the file using the adjusted request
	h.ServeHTTP(w, r)
}
