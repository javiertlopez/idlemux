package axiom

import (
	"context"
	"time"

	"github.com/javiertlopez/awesome/pkg/errorcodes"
	"github.com/javiertlopez/awesome/pkg/model"
	"github.com/javiertlopez/awesome/pkg/repository"

	guuid "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Collection keeps the collection name
const Collection = "videos"

// video model for mongodb
type video struct {
	ID          string    `bson:"_id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Duration    *float64  `bson:"duration,omitempty"`
	AssetID     string    `bson:"asset_id,omitempty"`
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}

// videos struct holds the logger and MongoDB client
type videos struct {
	db     string
	mongo  *mongo.Client
	logger *logrus.Logger
}

// NewVideoRepo method
func NewVideoRepo(
	l *logrus.Logger,
	db string,
	m *mongo.Client,
) repository.VideoRepo {
	return &videos{
		db:     db,
		mongo:  m,
		logger: l,
	}
}

// Create video creates a new ID, stores the video and returns the new object
func (v *videos) Create(ctx context.Context, anyVideo model.Video) (model.Video, error) {
	collection := v.mongo.Database(v.db).Collection(Collection)
	time := time.Now()

	uuid := guuid.New().String()

	insert := &video{
		ID:          uuid,
		Title:       anyVideo.Title,
		Description: anyVideo.Description,
		CreatedAt:   time,
		UpdatedAt:   time,
	}

	if anyVideo.Asset != nil {
		insert.AssetID = anyVideo.Asset.ID
	}

	_, err := collection.InsertOne(ctx, insert)
	if err != nil {
		v.logger.WithFields(log.Fields{
			"step": "collection.InsertOne",
			"func": "func (v *videos) Insert",
			"uuid": uuid,
		}).Error(err.Error())

		return model.Video{}, err
	}

	return insert.toModel(), nil
}

// GetByID retrieves a video with the ID
func (v *videos) GetByID(ctx context.Context, id string) (model.Video, error) {
	var response video

	collection := v.mongo.Database(v.db).Collection(Collection)

	filter := bson.M{"_id": id}

	err := collection.FindOne(ctx, filter).Decode(&response)

	if err != nil {
		v.logger.WithFields(log.Fields{
			"step": "collection.FindOne",
			"func": "func (v *videos) GetByID",
			"id":   id,
		}).Error(err.Error())

		if err == mongo.ErrNoDocuments {
			return model.Video{}, errorcodes.ErrVideoNotFound
		}

		return model.Video{}, err
	}

	return response.toModel(), nil
}

func (v video) toModel() model.Video {
	return model.Video{
		ID:          &v.ID,
		Title:       v.Title,
		Description: v.Description,
		Asset: &model.Asset{
			ID: v.AssetID,
		},
		Duration:  v.Duration,
		CreatedAt: v.CreatedAt.String(),
		UpdatedAt: v.UpdatedAt.String(),
	}
}
