package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type CommandLineFlags struct {
	addr           string
	fileServerAddr string
}

type Application struct {
	logErr  *log.Logger
	logInfo *log.Logger
	flags   *CommandLineFlags
}

func main() {
	app := &Application{
		logErr:  log.New(os.Stderr, "ERRORR\t", log.Lshortfile|log.Ltime|log.Ldate),
		logInfo: log.New(os.Stdout, "INFOO\t", log.Lshortfile|log.Ltime|log.Ldate),
		flags:   &CommandLineFlags{},
	}
	app.flags.getCommandLineFlags()

	server := &http.Server{
		Addr:     app.flags.addr,
		Handler:  app.routes(),
		ErrorLog: app.logErr,
	}

	db, err := pgx.Connect(context.Background(), "postgres://rinzler@/personalBlog")
	if err != nil {
		app.logErr.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())

	var id string
	var title string
	var content string
	var timestamp time.Time

	err = db.QueryRow(context.Background(), "select * from posts where id=1").Scan(&id, &title, &content, &timestamp)
	if err != nil {
		app.logErr.Println(err)
	}
	fmt.Printf("%v\n%v\n%v\n%v\n", id, title, content, timestamp)

	// Using my own error logger
	// logErr.Fatal(http.ListenAndServe(commandLineFlags.addr, mux))
	app.logInfo.Printf("Starting server on %s", app.flags.addr)
	app.logErr.Fatal(server.ListenAndServe())
}
