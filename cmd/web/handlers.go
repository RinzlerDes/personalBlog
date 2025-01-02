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

	if !formErrors.Valid() {
		app.renderPage(w, "postSearch.html", &ptd)
		return
	}

	// if targetPostIDString == "" {
	// 	ptd.FormErrors["id"] = models.FormErrorsState[models.EmptyFields]
	// 	// ptd.FormErrors.NotBlank(targetPostIDString)
	// 	ptd.InsertionErrorMessage = models.FormErrorsState[models.EmptyFields]
	// 	app.renderPage(w, "postSearch.html", &ptd)
	// 	return
	// }

	targetPostID, err := strconv.Atoi(targetPostIDString)
	// // targetPostID, err := strconv.Atoi(r.PathValue("id"))
	// if err != nil {
	// 	logErr.Println(err)
	// 	// ptd.IDIsNotNumber = true
	// 	ptd.InsertionErrorMessage = models.FormErrorsState[models.IDIsNotNumber]
	// 	app.renderPage(w, "postSearch.html", &ptd)
	// 	return
	// }
	//
	// logInfo.Println("should be post id:", targetPostID)
	// if targetPostID < 0 {
	// 	logInfo.Println("user's targetPostID is below 0")
	// 	// ptd.IDBelowZero = true
	// 	ptd.InsertionErrorMessage = models.FormErrorsState[models.IDBelowZero]
	// 	app.renderPage(w, "postSearch.html", &ptd)
	// 	return
	// }

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

	app.renderPage(w, "view.html", &ptd)

	// Either first render or error finding post or user input
	// app.renderPage(w, "postSearch.html", &ptd)
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

	// UNCOMMENT ALL OF THIS AND FIX

	// if ptd.Post.Title == "" {
	// 	// app.serverError(w, fmt.Errorf("Title nor content can be empty"))
	// 	// ptd.EmptyFields = true
	// 	ptd.FormErrors["title"] = models.FormErrorsState[models.EmptyFields]
	// 	ptd.InsertionErrorMessage = models.FormErrorsState[models.EmptyFields]
	// }
	//
	// if ptd.Post.Content == "" {
	// 	// app.serverError(w, fmt.Errorf("Title nor content can be empty"))
	// 	// ptd.EmptyFields = true
	// 	ptd.FormErrors["content"] = models.FormErrorsState[models.EmptyFields]
	// 	ptd.InsertionErrorMessage = models.FormErrorsState[models.EmptyFields]
	// }
	//
	// if len(ptd.FormErrors) > 0 {
	// 	app.renderPage(w, "insert.html", &ptd)
	// 	return
	// }

	err = app.postModel.Insert(&ptd.Post)
	if err != nil {
		logErr.Println(err)
		// ptd.PostInsertionError = true
		ptd.InsertionErrorMessage = models.FormErrorsState[models.PostInsertionError]
		app.renderPage(w, "insert.html", &ptd)
		return
	}
	logInfo.Println("post inserted")

	ptd.PostInserted = true

	app.renderPage(w, "insert.html", &ptd)
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
