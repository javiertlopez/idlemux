package router

import (
	"github.com/javiertlopez/awesome/controller"

	"github.com/gorilla/mux"
)

// setupVideoController setup the router with the event controller
func setupVideoController(router *mux.Router, video controller.VideoController) {
	router.HandleFunc("/videos", video.Create).Methods("POST")
	router.HandleFunc("/videos/{id}", video.GetByID).Methods("GET")
}
