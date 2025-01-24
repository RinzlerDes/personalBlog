package main

import (
	"errors"
	"fmt"
	"net/http"
	"personalBlog/internal/models"
	"strconv"
	"strings"
)

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	logInfo.Println("Home page")
	if r.URL.Path != "/" {
		logInfo.Printf("%s is an invalid path, serve http not found", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	posts, err := app.postModel.Latest(5)
	if err != nil {
		app.serverError(w, err)
		return
	}

	ptd := app.newPostTemplateData()
	ptd.Posts = posts
	app.renderPage(w, "home.html", &ptd)
}

func (app *Application) viewHandler(w http.ResponseWriter, r *http.Request) {
	logInfo.Println("In viewHandler")
	// Get ID from url
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 {
		str := fmt.Sprintf("%s\nPost %v does not exist\n", http.StatusText(http.StatusNotFound), id)
		http.Error(w, str, http.StatusNotFound)
		logErr.Printf("%v id=%v", err, id)
		return

	}

	logInfo.Printf("id: %d\n", id)

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
		// http.Error(w, err.Error(), http.StatusConflict)
		app.serverError(w, err)
		return
	}
	logInfo.Println("post found")

	sessionStr := app.sessionManager.PopString(r.Context(), "successfulInsert")
	logInfo.Printf("Popped string from session: %v\n", sessionStr)

	postTemplateData := app.newPostTemplateData()
	postTemplateData.Post = post
	postTemplateData.InitialInsertionMessage = sessionStr

	app.renderPage(w, "view.html", &postTemplateData)
}

func (app *Application) searchHandler(w http.ResponseWriter, r *http.Request) {
	logInfo.Println("Get Search Handler")
	postTemplateData := app.newPostTemplateData()
	app.renderPage(w, "postSearch.html", &postTemplateData)
}

func (app *Application) searchHandlerPost(w http.ResponseWriter, r *http.Request) {
	logInfo.Println("Post Search Handler")
	ptd := app.newPostTemplateData()

	// Limit amount of bytes able to be sent from a form
	// r.Body = http.MaxBytesReader(w, r.Body, 4)

	// Handle form input, feels weird to do this here
	err := r.ParseForm()
	if err != nil {
		logErr.Println(err)
		app.renderPage(w, "postSearch.html", &ptd)
		return
	}

	targetPostIDString := r.Form.Get("postID")

	formErrors := ptd.FormErrors

	formErrors.RunChecksForId(targetPostIDString, "id")

	if formErrors.NotValid() {
		app.renderPage(w, "postSearch.html", &ptd)
		return
	}

	targetPostID, _ := strconv.Atoi(targetPostIDString)

	// Get the post from DB using the ID
	ptd.Post, err = app.postModel.Get(uint(targetPostID))
	if err != nil {
		// No matching record error
		if errors.Is(err, models.ErrNoRecord) {
			logErr.Println("post not found in searchHandler: ", err)
			// ptd.PostNotFound = true
			ptd.InsertionErrorMessage = models.FormErrorsState[models.PostNotFound]

			app.renderPage(w, "postSearch.html", &ptd)
			return
		} else {
			// Internal server error
			app.serverError(w, err)
			return
		}
	}
	logInfo.Printf("got post %d from db\n", targetPostID)

	// app.renderPage(w, "view.html", &ptd)
	url := fmt.Sprintf("/posts/view/%d", ptd.Post.ID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (app *Application) insertHandler(w http.ResponseWriter, r *http.Request) {
	logInfo.Println("insertHandler")
	ptd := app.newPostTemplateData()
	app.renderPage(w, "insert.html", &ptd)
}

func (app *Application) insertHandlerPost(w http.ResponseWriter, r *http.Request) {
	logInfo.Println("insertHandlerPost")
	ptd := app.newPostTemplateData()

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "can't parse form", http.StatusBadRequest)
		return
	}

	ptd.Post.Title = strings.TrimSpace(r.Form.Get("title"))
	ptd.Post.Content = strings.TrimSpace(r.Form.Get("content"))
	logInfo.Println(ptd.Post)

	formErrors := ptd.FormErrors
	formErrors.RunChecksForTitle(ptd.Post.Title, "title")
	formErrors.RunChecksForContent(ptd.Post.Content, "content")

	if formErrors.NotValid() {
		app.renderPage(w, "insert.html", &ptd)
		return
	}

	id, err := app.postModel.Insert(&ptd.Post)
	if err != nil {
		logErr.Println(err)
		// ptd.PostInsertionError = true
		ptd.InsertionErrorMessage = models.FormErrorsState[models.PostInsertionError]
		app.renderPage(w, "insert.html", &ptd)
		return
	}
	logInfo.Println("post inserted")

	ptd.PostInserted = true

	app.sessionManager.Put(r.Context(), "successfulInsert", "Post Inserted Successfully")

	// Redirect to new post
	url := fmt.Sprintf("/posts/view/%d", id)
	http.Redirect(w, r, url, http.StatusSeeOther)
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
