package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our applicaiton routes.
func (app *application) routes() http.Handler {
	// mux := http.NewServeMux()

	router := httprouter.New()

	// set the handler for 404 and 405
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Update the pattern for the route for the static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Create a new middleware chain containing the middleware specific to our dynmic routes
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// Create the routes using the appropriate methods, patterns and handlers.

	// router.HandlerFunc(http.MethodGet, "/", app.home)
	// router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	// router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	// router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}
