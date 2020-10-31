package main

import (
	"context"
	"time"

	guuid "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	logger *logrus.Logger
	mongo  *mongo.Client
}

// NewVideoService creates new an Videos service object
func NewVideoService(
	l *logrus.Logger,
	m *mongo.Client,
) Videos {
	return &videos{
		logger: l,
		mongo:  m,
	}
}

// Insert video creates a new ID, stores the video and returns the new object
func (v *videos) Insert(ctx context.Context, anyVideo *Video) (*Video, error) {
	collection := v.mongo.Database("awesome").Collection("videos")
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

		return nil, err
	}

	return insert.toModel(), nil
}

// GetByID retrieves a video with the ID
func (v *videos) GetByID(ctx context.Context, id string) (*Video, error) {
	var response video

	collection := v.mongo.Database("awesome").Collection("videos")

	filter := bson.M{"_id": id}

	err := collection.FindOne(ctx, filter).Decode(&response)

	if err != nil {
		v.logger.WithFields(log.Fields{
			"step": "collection.FindOne",
			"func": "func (v *videos) GetByID",
			"id":   id,
		}).Error(err.Error())

		if err == mongo.ErrNoDocuments {
			return nil, ErrVideoNotFound
		}

		return nil, err
	}

	anyVideo := response.toModel()

	return anyVideo, nil
}

func (v video) toModel() *Video {
	return &Video{
		ID:          &v.ID,
		Title:       v.Title,
		Description: v.Description,
		Asset: &Asset{
			ID: v.AssetID,
		},
		Duration:  v.Duration,
		CreatedAt: v.CreatedAt.String(),
		UpdatedAt: v.UpdatedAt.String(),
	}
}
