package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// addAppHandler adds the handler to the mux router
func (a *App) addAppHandler(r *mux.Router) {
	r.HandleFunc("/app/hello", a.HelloHandler).Methods("GET")
}

// HelloHandler handler
func (a *App) HelloHandler(w http.ResponseWriter, r *http.Request) {
	JSONResponse(
		w,
		http.StatusOK,
		Response{
			Message: "Hello World!",
			Status:  http.StatusOK,
		},
	)
}
