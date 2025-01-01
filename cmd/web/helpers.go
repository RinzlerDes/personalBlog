package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"personalBlog/internal/models"
	"runtime/debug"
)

// func (app *Application) serverError(w http.ResponseWriter, err error) {
func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// app.logErr.Output(2, trace)
	logErr.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (flags *CommandLineFlags) getCommandLineFlags() {
	flag.StringVar(&flags.addr, "addr", "localhost:8080", "HTTP network address")
	flag.StringVar(&flags.fileServerAddr, "fileServerAddr", "./ui/static", "Path to static assets")
	flag.Parse()
}

func (app *Application) renderPage(w http.ResponseWriter, templateName string, ptd *models.PostTemplateData) {
	template, exists := app.parsedTemplatesCache[templateName]
	if !exists {
		app.serverError(w, fmt.Errorf("the page for %s does not exist", templateName))
		return
	}

	var buf bytes.Buffer

	err := template.ExecuteTemplate(&buf, "base", ptd)
	if err != nil {
		app.serverError(w, err)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		logErr.Println("error writing buf")
		app.serverError(w, err)
	}
}

func (app *Application) newPostTemplateData() models.PostTemplateData {
	return models.NewPostTemplateData()
	// return models.PostTemplateData{
	// 	CurrentYear: time.Now().Year(),
	// }
}
