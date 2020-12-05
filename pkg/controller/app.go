package controller

import (
	"net/http"
)

// AppController handles the HTTP requests
type AppController interface {
	Healthz(w http.ResponseWriter, r *http.Request)
	Statusz(w http.ResponseWriter, r *http.Request)
}

// appController struct holds the usecase
type appController struct {
	commit  string
	version string
}

// NewAppController returns an EventController
func NewAppController(commit string, version string) AppController {
	return &appController{
		commit,
		version,
	}
}

// Healthz controller
func (ac *appController) Healthz(w http.ResponseWriter, r *http.Request) {
	JSONResponse(
		w,
		http.StatusOK,
		Response{
			Message: "Hello World!",
			Status:  http.StatusOK,
		},
	)
}

// Statusz controller
func (ac *appController) Statusz(w http.ResponseWriter, r *http.Request) {
	JSONResponse(
		w,
		http.StatusOK,
		map[string]interface{}{
			"commit":  ac.commit,
			"version": ac.version,
		},
	)
}
