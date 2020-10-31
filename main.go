package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// App holds the Video service (seems okay)
type App struct {
	logger *logrus.Logger
	videos Videos
}

func main() {
	mongoString := os.Getenv("MONGO_STRING")

	// Create a logrus logger and set up the output format as JSON
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	// Connect to Mongo Atlas
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(
		mongoString,
	))
	if err != nil {
		logger.Fatal(err)
	}

	// Create an app instance
	awesome := &App{
		logger: logger,
		videos: NewVideoService(
			logger,
			mongoClient,
		),
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
