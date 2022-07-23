package awesome

import (
	"context"
	"time"

	"github.com/gorilla/mux"
	muxgo "github.com/muxinc/mux-go/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/javiertlopez/awesome/controller"
	"github.com/javiertlopez/awesome/repository/axiom"
	"github.com/javiertlopez/awesome/repository/muxinc"
	"github.com/javiertlopez/awesome/router"
	"github.com/javiertlopez/awesome/usecase"
)

const (
	Database     = "delivery"       // Database keeps the database name
	mongoTimeout = 15 * time.Second // mongotimeout
)

// App holds the handler, and logger
type App struct {
	logger *logrus.Logger
	router *mux.Router
}

// AppConfig struct with configuration variables
type AppConfig struct {
	Commit         string
	Version        string
	MongoURI       string
	MuxTokenID     string
	MuxTokenSecret string
	MuxKeyID       string
	MuxKeySecret   string
	Test           bool
}

// New returns an App
func New(config AppConfig, logger *logrus.Logger) App {
	// Set client options
	clientOptions := options.Client().ApplyURI(config.MongoURI)

	// Context with timeout for establish connection with Mongo Atlas
	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
	defer cancel()

	// Connect to Mongo Atlas
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal(err)
	}
	db := client.Database(Database)

	// Init mux repository
	assets := muxinc.New(
		logger,
		muxgo.NewAPIClient(
			muxgo.NewConfiguration(
				muxgo.WithBasicAuth(config.MuxTokenID, config.MuxTokenSecret),
			),
		),
		muxinc.Config{
			KeyID:     config.MuxKeyID,
			KeySecret: config.MuxKeySecret,
			Test:      config.Test,
		},
	)

	// Init axiom repository
	videos := axiom.New(logger, db)

	// Init delivery usecase
	delivery := usecase.Delivery(assets, videos)

	// Init ingestion usecase
	ingestion := usecase.Ingestion(assets, videos)

	// Init controller
	controller := controller.New(config.Commit, config.Version, delivery, ingestion)

	// Setup router
	router := router.New(controller)

	return App{
		logger: logger,
		router: router,
	}
}

// Router returns the *mux.Router
func (a *App) Router() *mux.Router {
	return a.router
}
