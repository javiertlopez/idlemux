package router

import (
	"github.com/javiertlopez/awesome/controller"

	"github.com/gorilla/mux"
)

// setupAppController setup the router with the app controller
func setupAppController(router *mux.Router, app controller.AppController) {
	router.HandleFunc("/app/healthz", app.Healthz).Methods("GET")
	router.HandleFunc("/app/statusz", app.Statusz).Methods("GET")
}
