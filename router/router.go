package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Controller handles the HTTP requests
type Controller interface {
	Healthz(w http.ResponseWriter, r *http.Request)
	Statusz(w http.ResponseWriter, r *http.Request)

	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

// New returns a *mux.Router
func New(
	controller Controller,
) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/app/healthz", controller.Healthz).Methods("GET")
	router.HandleFunc("/app/statusz", controller.Statusz).Methods("GET")

	router.HandleFunc("/videos", controller.Create).Methods("POST")
	router.HandleFunc("/videos/{id}", controller.GetByID).Methods("GET")
	router.HandleFunc("/videos", controller.List).Methods("GET")

	return router
}
