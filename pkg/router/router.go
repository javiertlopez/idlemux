package router

import (
	"github.com/gorilla/mux"
	"github.com/javiertlopez/awesome/pkg/controller"
)

// New returns a *mux.Router
func New(video controller.VideoController) *mux.Router {
	router := mux.NewRouter()

	setupVideoController(router, video)

	return router
}
