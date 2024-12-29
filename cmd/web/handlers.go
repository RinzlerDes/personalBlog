package main

import (
	"errors"
	"fmt"
	"net/http"
	"personalBlog/internal/models"
	"strconv"
)

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	logInfo.Println("Home page")
	if r.URL.Path != "/" {
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
		http.Error(w, err.Error(), http.StatusConflict)
		app.serverError(w, err)
		return
	}

	postTemplateData := app.newPostTemplateData()
	postTemplateData.Post = post

	app.renderPage(w, "view.html", &postTemplateData)
}

// func (app *Application) viewHandler(w http.ResponseWriter, r *http.Request) {
// 	// Get ID from url
// 	id, err := strconv.Atoi(r.URL.Query().Get("id"))
// 	if err != nil || id < 0 {
// 		str := fmt.Sprintf("%s\nPost %v does not exist\n", http.StatusText(http.StatusNotFound), id)
// 		http.Error(w, str, http.StatusNotFound)
// 		logErr.Printf("%v id=%v", err, id)
// 		return
// 	}
//
// 	// Get the post from DB using the ID
// 	post, err := app.postModel.Get(uint(id))
// 	if err != nil {
// 		// No matching record error
// 		if errors.Is(err, models.ErrNoRecord) {
// 			logErr.Println("post not found: ", err)
// 			http.Error(w, err.Error(), http.StatusNotFound)
// 			return
// 		}
//
// 		// Internal server error
// 		http.Error(w, err.Error(), http.StatusConflict)
// 		app.serverError(w, err)
// 		return
// 	}
//
// 	postTemplateData := app.newPostTemplateData()
// 	postTemplateData.Post = post
//
// 	app.renderPage(w, "view.html", &postTemplateData)
// }

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
	logInfo.Println("searchHandler")
	postTemplateData := app.newPostTemplateData()
	app.renderPage(w, "postSearch.html", &postTemplateData)
}

func (app *Application) searchHandlerProcessForm(w http.ResponseWriter, r *http.Request) {
	logInfo.Println("Search Page")
	ptd := app.newPostTemplateData()

	// Handle form input, feels weird to do this here
	err := r.ParseForm()
	if err != nil {
		logErr.Println(err)
		app.renderPage(w, "postSearch.html", nil)
		return
	}

	targetPostIDString := r.Form.Get("postID")
	if targetPostIDString == "" {
		ptd.InsertionErrorMessage = models.InsertionErrorsState[models.EmptyFields]
		app.renderPage(w, "postSearch.html", &ptd)
		return
	}

	targetPostID, err := strconv.Atoi(targetPostIDString)
	// targetPostID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		logErr.Println(err)
		// ptd.IDIsNotNumber = true
		ptd.InsertionErrorMessage = models.InsertionErrorsState[models.IDIsNotNumber]
		app.renderPage(w, "postSearch.html", &ptd)
		return
	}

	logInfo.Println("should be post id:", targetPostID)
	if targetPostID < 0 {
		logInfo.Println("user's targetPostID is below 0")
		// ptd.IDBelowZero = true
		ptd.InsertionErrorMessage = models.InsertionErrorsState[models.IDBelowZero]
		app.renderPage(w, "postSearch.html", &ptd)
		return
	}

	// Get the post from DB using the ID
	ptd.Post, err = app.postModel.Get(uint(targetPostID))
	if err != nil {
		// No matching record error
		if errors.Is(err, models.ErrNoRecord) {
			logErr.Println("post not found in searchHandler: ", err)
			// ptd.PostNotFound = true
			ptd.InsertionErrorMessage = models.InsertionErrorsState[models.PostNotFound]

			app.renderPage(w, "postSearch.html", &ptd)
			return
		} else {
			// Internal server error
			app.serverError(w, err)
			return
		}
	}
	logInfo.Printf("got post %d from db\n", targetPostID)

	app.renderPage(w, "view.html", &ptd)

	// Either first render or error finding post or user input
	// app.renderPage(w, "postSearch.html", &ptd)
}

func (app *Application) insertHandler(w http.ResponseWriter, r *http.Request) {
	logInfo.Println("Insert Page")
	ptd := app.newPostTemplateData()

	if r.Method == "POST" {
		app.insertPostLogic(w, r, &ptd)
	} else {
		app.renderPage(w, "insert.html", &ptd)
	}
}

func (app *Application) insertPostLogic(w http.ResponseWriter, r *http.Request, ptd *models.PostTemplateData) {
	r.ParseForm()
	title := r.Form.Get("title")
	content := r.Form.Get("content")

	if title == "" || content == "" {
		// app.serverError(w, fmt.Errorf("Title nor content can be empty"))
		// ptd.EmptyFields = true
		ptd.InsertionErrorMessage = models.InsertionErrorsState[models.EmptyFields]
		app.renderPage(w, "insert.html", ptd)
		return
	}

	post := models.Post{
		Title:   title,
		Content: content,
	}

	err := app.postModel.Insert(&post)
	if err != nil {
		logErr.Println(err)
		// ptd.PostInsertionError = true
		ptd.InsertionErrorMessage = models.InsertionErrorsState[models.PostInsertionError]
		app.renderPage(w, "insert.html", ptd)
		return
	}
	logInfo.Println("post inserted")

	// ptd.PostInserted = true
	app.renderPage(w, "insert.html", ptd)
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
