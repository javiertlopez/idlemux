package mongodb

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/javiertlopez/awesome/errorcodes"
	"github.com/javiertlopez/awesome/model"
)

// Collection keeps the collection name
const Collection = "videos"

// video model for mongodb
type video struct {
	ID          string    `bson:"_id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Duration    float64   `bson:"duration,omitempty"`
	AssetID     string    `bson:"asset_id,omitempty"`
	MasterID    string    `bson:"master_id,omitempty"`
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}

// Create video creates a new ID, stores the video and returns the new object
func (db *DB) Create(ctx context.Context, anyVideo model.Video) (model.Video, error) {
	collection := db.mongo.Collection(Collection)
	time := time.Now()

	id := uuid.New().String()

	insert := &video{
		ID:          id,
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
		db.logger.WithFields(logrus.Fields{
			"step": "collection.InsertOne",
			"func": "func (v *videos) Insert",
			"id":   id,
		}).Error(err.Error())

		return model.Video{}, err
	}

	return insert.toModel(), nil
}

// GetByID retrieves a video with the ID
func (db *DB) GetByID(ctx context.Context, id string) (model.Video, error) {
	var response video

	collection := db.mongo.Collection(Collection)

	filter := bson.D{{Key: "_id", Value: id}}

	err := collection.FindOne(ctx, filter).Decode(&response)

	if err != nil {
		db.logger.WithFields(logrus.Fields{
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
		ID:          v.ID,
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
