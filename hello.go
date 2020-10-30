package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// addAppHandler adds the handler to the mux router
func addAppHandler(r *mux.Router) {
	r.HandleFunc("/app/hello", HelloHandler).Methods("GET")
}

// HelloHandler handler
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	JSONResponse(
		w,
		http.StatusOK,
		Response{
			Message: "Hello World!",
			Status:  http.StatusOK,
		},
	)
}
