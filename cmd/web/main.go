package main

import (
	"context"
	_ "fmt"
	"net/http"
	_ "time"

	"personalBlog/internal/loggers"
	"personalBlog/internal/models"

	"github.com/jackc/pgx/v5"
)

type CommandLineFlags struct {
	addr           string
	fileServerAddr string
}

type Application struct {
	flags *CommandLineFlags
	posts *models.PostModel
}

// Create loggers
var logErr = loggers.LogErr
var logInfo = loggers.LogInfo

func main() {
	// // Create loggers
	// logErr := loggers.LogErr
	// logInfo := loggers.LogInfo

	// Open database connection
	db, err := openDB("postgres://rinzler@/personalBlog")
	if err != nil {
		logErr.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())

	app := &Application{
		flags: &CommandLineFlags{},
		posts: &models.PostModel{DB: db},
	}

	app.flags.getCommandLineFlags()

	server := &http.Server{
		Addr:    app.flags.addr,
		Handler: app.routes(),
	}

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

// func testInsert(motelPost *models.Post) {
//     	newPost := models.Post{
// 		Title:   "First insert from within go",
// 		Content: "More like second try",
// 	}
//
// 	err = app.posts.Insert(&newPost)
// 	if err != nil {
// 		logErr.Println(err)
// 	}
// 	logInfo.Printf("new post\nID        %d\nTitle       %s\nContent     %s\nCreated     %v", newPost.ID, newPost.Title, newPost.Content, newPost.Created)
// }
