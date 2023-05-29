package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sxc/snippetbox/internal/models"

	"github.com/alexedwards/scs/v2"

	"github.com/alexedwards/scs/mysqlstore"

	"github.com/go-playground/form/v4"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// flag will be stored in the addr variable at run time.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "web:password@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	// Use log.New() to creat a logger for writing information message.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ldate|log.Ltime)
	// Use log.New() to creat a logger for writing error messages.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create a connection pool
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new form decoder.
	formDecoder := form.NewDecoder()

	// Initialize a new session manager.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Initialize a new instance of the application struct.
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// tls.Config struct to hold the non-default TLS settings
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		// MinVersion:       tls.VersionTLS12,
		// MaxVersion:       tls.VersionTLS12,
	}

	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,

		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// srv := &http.Server{
	// 	Addr:     *addr,
	// 	ErrorLog: errorLog,
	// 	Handler:  app.routes(),
	// }

	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.
	infoLog.Printf("Starting server on %s", *addr)
	// err = srv.ListenAndServe()
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// openDB() returns a sql.DB connection pool.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
