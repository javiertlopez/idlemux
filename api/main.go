package main

import (
	"context"
	"net/http"
	"os"
	"time"

	awesome "github.com/javiertlopez/awesome/pkg"
	"github.com/javiertlopez/awesome/pkg/asset"
	"github.com/javiertlopez/awesome/pkg/video"
	muxgo "github.com/muxinc/mux-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// App holds the Video service (seems okay)
type App struct {
	logger *logrus.Logger
	assets awesome.Assets
	videos awesome.Videos
}

func main() {
	// Environment variables
	mongoString := os.Getenv("MONGO_STRING")
	muxTokenID := os.Getenv("MUX_TOKEN_ID")
	muxTokenSecret := os.Getenv("MUX_TOKEN_SECRET")
	muxKeyID := os.Getenv("MUX_KEY_ID")
	muxKeySecret := os.Getenv("MUX_KEY_SECRET")

	// Create a logrus logger and set up the output format as JSON
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	// Context with timeout for establish connection with Mongo Atlas
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Connect to Mongo Atlas
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(
		mongoString,
	))
	if err != nil {
		logger.Fatal(err)
	}

	// Mux.com API Client Initialization
	muxClient := muxgo.NewAPIClient(
		muxgo.NewConfiguration(
			muxgo.WithBasicAuth(muxTokenID, muxTokenSecret),
		),
	)

	// Create an app instance
	awesomeApp := &App{
		logger: logger,
		assets: asset.NewAssetService(
			logger,
			muxClient,
			muxKeyID,
			muxKeySecret,
		),
		videos: video.NewVideoService(
			logger,
			mongoClient,
		),
	}

	// Create a Gorilla Mux router
	router := awesomeApp.Router()

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
