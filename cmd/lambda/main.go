package main

import (
	"context"
	"os"

	awesome "github.com/javiertlopez/awesome/pkg"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/sirupsen/logrus"
)

var (
	application awesome.App
	commit      string
	version     string
	adapter     *gorillamux.GorillaMuxAdapter
)

func init() {
	// Environment variables
	mongoDB := os.Getenv("MONGO_DB")
	mongoString := os.Getenv("MONGO_STRING")
	muxTokenID := os.Getenv("MUX_TOKEN_ID")
	muxTokenSecret := os.Getenv("MUX_TOKEN_SECRET")
	muxKeyID := os.Getenv("MUX_KEY_ID")
	muxKeySecret := os.Getenv("MUX_KEY_SECRET")

	// Create a logrus logger and set up the output format as JSON
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	application = awesome.New(
		awesome.AppConfig{
			Commit:         commit,
			Version:        version,
			MongoDB:        mongoDB,
			MongoURI:       mongoString,
			MuxTokenID:     muxTokenID,
			MuxTokenSecret: muxTokenSecret,
			MuxKeyID:       muxKeyID,
			MuxKeySecret:   muxKeySecret,
		},
		logger,
	)

	// Create a Gorilla Mux router
	adapter = gorillamux.New(application.Router())
}

// Handler for lamba adapter
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adapter.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
