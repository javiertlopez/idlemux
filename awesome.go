package awesome

import (
	"context"
	"time"

	"github.com/javiertlopez/awesome/controller"
	"github.com/javiertlopez/awesome/repository/axiom"
	"github.com/javiertlopez/awesome/repository/muxinc"
	"github.com/javiertlopez/awesome/router"
	"github.com/javiertlopez/awesome/usecase"

	"github.com/gorilla/mux"
	muxgo "github.com/muxinc/mux-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoTimeout = 15 * time.Second

// App holds the handler, and logger
type App struct {
	logger *logrus.Logger
	router *mux.Router
	config AppConfig
}

// AppConfig struct with configuration variables
type AppConfig struct {
	Commit         string
	Version        string
	MongoDB        string
	MongoURI       string
	MuxTokenID     string
	MuxTokenSecret string
	MuxKeyID       string
	MuxKeySecret   string
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

	// Init mux repository
	assetsRepo := muxinc.NewAssetRepo(
		logger,
		muxgo.NewAPIClient(
			muxgo.NewConfiguration(
				muxgo.WithBasicAuth(config.MuxTokenID, config.MuxTokenSecret),
			),
		),
		config.MuxKeyID,
		config.MuxKeySecret,
	)

	// Init axiom repository
	videosRepo := axiom.NewVideoRepo(logger, config.MongoDB, client)

	// Init usecase
	videos := usecase.NewVideoUseCase(assetsRepo, videosRepo)

	// Init appController
	appController := controller.NewAppController(config.Commit, config.Version)

	// Init videoController
	videoController := controller.NewVideoController(videos)

	// Setup router
	router := router.New(appController, videoController)

	return App{
		logger,
		router,
		config,
	}
}

// Router returns the *mux.Router
func (a *App) Router() *mux.Router {
	return a.router
}
