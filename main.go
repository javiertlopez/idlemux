package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	// Create a logrus logger and set up the output format as JSON
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	// Create a Gorilla Mux router
	router := mux.NewRouter()

	// Pass the router to set up handlers
	addAppHandler(router)
	addVideoHandler(router)

	// Create a Server instance with the router
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      router,
	}

	// Start the server
	logger.Fatal(srv.ListenAndServe())
}
