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

	postTemplateData := models.PostTemplateData{
		Post: post,
	}

	err = t.ExecuteTemplate(w, "base", postTemplateData)
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
	postTemplateData := models.PostTemplateData{}

	// Handle form input, feels weird to do this here
	if r.Method == "POST" {
		r.ParseForm()

		postErrors := false

		targetPostID, err := strconv.Atoi(r.Form.Get("postID"))
		if err != nil {
			logErr.Println(err)
			postTemplateData.IDIsNotNumber = true
			postErrors = true
		}
		logInfo.Println("should be post id:", targetPostID)
		if targetPostID < 0 {
			logInfo.Println("user's targetPostID is below 0")
			postTemplateData.IDBelowZero = true
			postErrors = true
		}

		if !postErrors {
			recordFound := true

			// Get the post from DB using the ID
			postTemplateData.Post, err = app.postModel.Get(uint(targetPostID))
			if err != nil {
				// No matching record error
				if errors.Is(err, models.ErrNoRecord) {
					logErr.Println("post not found in searchHandler: ", err)
					postTemplateData.PostNotFound = true
					postErrors = true
					recordFound = false
				} else {
					// Internal server error
					app.serverError(w, err)
					return
				}
			}

			if recordFound { // Try using cookies to redirect with already fetched struct
				// All good, redirect to view
				// http.Redirect(w http.ResponseWriter, r *http.Request, url string, code int)
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

				err = t.ExecuteTemplate(w, "base", postTemplateData)
				if err != nil {
					app.serverError(w, err)
					return
				}
				return
			}
		}
	}

	// Either first render or error finding post or user input
	app.RenderSearchTemplate(w, postTemplateData)
}

func (app *Application) RenderSearchTemplate(w http.ResponseWriter, postTemplateData models.PostTemplateData) {
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

	err = t.ExecuteTemplate(w, "base", postTemplateData)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *Application) insertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		app.insertPostLogic(w, r)
	} else {
		app.renderInsertHTML(w, models.PostTemplateData{})
	}
}

func (app *Application) insertPostLogic(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	content := r.Form.Get("content")

	post := models.Post{
		Title:   title,
		Content: content,
	}

	ptd := models.PostTemplateData{}

	err := app.postModel.Insert(&post)
	if err != nil {
		logErr.Println(err)
		ptd.PostInsertionError = true
		app.renderInsertHTML(w, ptd)
		return
	}
	logInfo.Println("post inserted")

	ptd.PostInserted = true
	app.renderInsertHTML(w, ptd)
}

func (app *Application) renderInsertHTML(w http.ResponseWriter, ptd models.PostTemplateData) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/insert.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = t.ExecuteTemplate(w, "base", ptd)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *Application) fileServerHandler(w http.ResponseWriter, r *http.Request, h http.Handler) {
	// orig := r.URL.Path
	strippedPath := r.URL.Path[len("/static/"):] // Get the path after stripping
	// logInfo.Printf("Before: %s\nFile path after StripPrefix: %s", orig, strippedPath)

	// Adjust the request URL path to match the file server's expectation
	r.URL.Path = strippedPath // Set the adjusted path for the file server

	// Serve the file using the adjusted request
	h.ServeHTTP(w, r)
}
