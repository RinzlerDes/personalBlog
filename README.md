A personal blog web app written in GO. The web app consists of 3 pages at the moment, a home, search, view and insert page. The home page displays the most recent posts. The search and insert page take in user input via a form and redirects to a view page if its request is successful. User and session data is stored in PostgreSQL relational database. The website pages themselves are comprised of HTML and GO's HTML/Template package to populate the webpage with data at runtime.

Some notable libraries used
net/http
This package in the standard module is used to set up the server, establish routes to the website and is used in handlers to modify certain HTTP properties.
github.com/jackc/pgx/v5/pgxpool
Used to setup connection pool with my PostgeSQL database.

Command to run: go run ./cmd/web/
Flags: 
-help
-addr custom address where program will be hosted, default in localhost:8080
-fileServerAddr custom address where files used on website will pulled from, default is ./ui/static

