package main

import (
	"context"
	"net/http"
	"personalBlog/internal/loggers"
	"personalBlog/internal/models"
	"text/template"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommandLineFlags struct {
	addr           string
	fileServerAddr string
}

type Application struct {
	flags                *CommandLineFlags
	postModel            *models.PostModel
	parsedTemplatesCache map[string]*template.Template
	sessionManager       *scs.SessionManager
}

// Create loggers
var (
	logErr  = loggers.LogErr
	logInfo = loggers.LogInfo
)

func main() {
	// Open database connection
	dbPool, err := openDB("postgres://rinzler@/personalBlog")
	if err != nil {
		logErr.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	// Parse and cache templates
	templates, err := newTemplateCache()
	if err != nil {
		logErr.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(dbPool)
	sessionManager.Lifetime = 1 * time.Minute

	app := &Application{
		flags:                &CommandLineFlags{},
		postModel:            &models.PostModel{DBPool: dbPool},
		parsedTemplatesCache: templates,
		sessionManager:       sessionManager,
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

func openDB(dsn string) (*pgxpool.Pool, error) {
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

func (app *Application) restInsert(newPost *models.Post) {
	_, err := app.postModel.Insert(newPost)
	if err != nil {
		logErr.Println(err)
		return
	}
	logInfo.Printf("new post\nID        %d\nTitle       %s\nContent     %s\nCreated     %v", newPost.ID, newPost.Title, newPost.Content, newPost.Created)
}
