package router

import (
	"github.com/gorilla/mux"
	"github.com/javiertlopez/awesome/pkg/controller"
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
