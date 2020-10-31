package main

import "github.com/gorilla/mux"

// Router creates a *mux.Router, setup handlers and return the router
func (a *App) Router() *mux.Router {
	// Create a Gorilla Mux router
	router := mux.NewRouter()

	// Pass the router to set up handlers
	a.addAppHandler(router)
	a.addVideoHandler(router)

	return router
}
