package awesome

import (
	"time"

	"github.com/gorilla/mux"
	muxgo "github.com/muxinc/mux-go/v5"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

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

	// Connect to Mongo Atlas
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		logger.WithError(err).Error(err.Error())
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
	delivery := usecase.Delivery(assets, videos, logger)

	// Init ingestion usecase
	ingestion := usecase.Ingestion(assets, videos, logger)

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
