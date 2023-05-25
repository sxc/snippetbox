package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our applicaiton routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/create", app.snippetCreate)
	mux.HandleFunc("/snippet/view", app.snippetView)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)
}
