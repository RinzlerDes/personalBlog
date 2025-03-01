package main

import (
	"fmt"
	"net/http"
	"personalBlog/internal/loggers"
	"personalBlog/internal/models"
	"text/template"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
)

type CommandLineFlags struct {
	addr           string
	fileServerAddr string
}

type Application struct {
	flags                CommandLineFlags
	postModel            *models.PostModel
	userModel            *models.UserModel
	parsedTemplatesCache map[string]*template.Template
	sessionManager       *scs.SessionManager
}

// Create loggers
var (
	logErr  = loggers.LogErr
	logInfo = loggers.LogInfo
)

func main() {
	commandLineFlags := getCommandLineFlags()
	fmt.Println(commandLineFlags)

	// Open database connection
	dbPool, err := openDBPool("postgres://rinzler@/personalBlog")
	if err != nil {
		logErr.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	// Parse and cache templates
	templates, err := newTemplateCache()
	if err != nil {
		logErr.Fatal(err)
	}

	// Create session manager
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(dbPool)
	sessionManager.Lifetime = 1 * time.Minute
	sessionManager.Cookie.Secure = true

	// Create app data
	app := &Application{
		flags:                commandLineFlags,
		postModel:            &models.PostModel{DBPool: dbPool},
		userModel:            &models.UserModel{DBPool: dbPool},
		parsedTemplatesCache: templates,
		sessionManager:       sessionManager,
	}

	// Create server
	server := &http.Server{
		Addr:    app.flags.addr,
		Handler: app.routes(),
	}

	logInfo.Printf("Starting server on %s", app.flags.addr)
	// Run server
	// logErr.Fatal(server.ListenAndServe())
	logErr.Fatal(server.ListenAndServeTLS("./transportLayerSecurity/cert.pem", "./transportLayerSecurity/key.pem"))
}
