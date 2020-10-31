package main

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// App holds the Video service (seems okay)
type App struct {
	logger *logrus.Logger
	videos Videos
}

func main() {
	// Create a logrus logger and set up the output format as JSON
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	// Create an app instance
	awesome := &App{
		logger: logger,
	}

	// Create a Gorilla Mux router
	router := awesome.Router()

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
