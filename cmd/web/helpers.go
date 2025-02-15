package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"personalBlog/internal/models"
	"runtime/debug"

	"github.com/jackc/pgx/v5/pgxpool"
)

func openDBPool(dsn string) (*pgxpool.Pool, error) {
	// db, err := pgx.Connect(context.Background(), dsn)
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}

// func (app *Application) serverError(w http.ResponseWriter, err error) {
func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// app.logErr.Output(2, trace)
	logErr.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Parse command line arguments and store them in flags struct
// func (flags *CommandLineFlags) getCommandLineFlags() {
func getCommandLineFlags() CommandLineFlags {
	flags := CommandLineFlags{}
	flag.StringVar(&flags.addr, "addr", "localhost:8080", "HTTP network address")
	flag.StringVar(&flags.fileServerAddr, "fileServerAddr", "./ui/static", "Path to static assets")
	flag.Parse()
	return flags
}

func (app *Application) renderPage(w http.ResponseWriter, templateName string, ptd models.TemplateData) {
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
