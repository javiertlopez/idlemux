package router

import (
	"github.com/javiertlopez/awesome/controller"

	"github.com/gorilla/mux"
)

// New returns a *mux.Router
func New(
	app controller.AppController,
	video controller.VideoController,
) *mux.Router {
	router := mux.NewRouter()

	setupAppController(router, app)
	setupVideoController(router, video)

	return router
}
