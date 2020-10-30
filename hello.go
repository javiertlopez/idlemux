package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Response struct
type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// addAppHandler adds the handler to the mux router
func addAppHandler(r *mux.Router) {
	r.HandleFunc("/app/hello", HelloHandler).Methods("GET")
}

// HelloHandler handler
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "Hello World!",
		Status:  http.StatusOK,
	}

	output, _ := json.Marshal(response)

	// Set the content type to json for browsers
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
