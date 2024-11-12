package main

import (
	"context"
	"net/http"
	"personalBlog/internal/loggers"
	"personalBlog/internal/models"
	"text/template"
	_ "time"

	"github.com/jackc/pgx/v5"
)

type CommandLineFlags struct {
	addr           string
	fileServerAddr string
}

type Application struct {
	flags                *CommandLineFlags
	postModel            *models.PostModel
	parsedTemplatesCache map[string]*template.Template
}

// Create loggers
var (
	logErr  = loggers.LogErr
	logInfo = loggers.LogInfo
)

func main() {
	// Open database connection
	db, err := openDB("postgres://rinzler@/personalBlog")
	if err != nil {
		logErr.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())

	// Parse and cache templates
	templates, err := newTemplateCache()
	if err != nil {
		logErr.Fatal(err)
	}

	app := &Application{
		flags:                &CommandLineFlags{},
		postModel:            &models.PostModel{DB: db},
		parsedTemplatesCache: templates,
	}

	app.flags.getCommandLineFlags()

	server := &http.Server{
		Addr:    app.flags.addr,
		Handler: app.routes(),
	}

	// -------------------------------------------------------------------------------------

	// -------------------------------------------------------------------------------------

	logInfo.Printf("Starting server on %s", app.flags.addr)
	logErr.Fatal(server.ListenAndServe())
}

func openDB(dsn string) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *Application) testInsert(newPost *models.Post) {
	err := app.postModel.Insert(newPost)
	if err != nil {
		logErr.Println(err)
		return
	}
	logInfo.Printf("new post\nID        %d\nTitle       %s\nContent     %s\nCreated     %v", newPost.ID, newPost.Title, newPost.Content, newPost.Created)
}
