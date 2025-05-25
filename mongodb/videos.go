package mongodb

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/javiertlopez/idlemux/errorcodes"
	"github.com/javiertlopez/idlemux/model"
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
		db.logger.WithError(err).Error("error inserting video into collection")

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
		db.logger.WithError(err).Error("error getting video by ID")

		if err == mongo.ErrNoDocuments {
			return model.Video{}, errorcodes.ErrVideoNotFound
		}

		return model.Video{}, err
	}

	return response.toModel(), nil
}

// List returns paginated videos from the collection using page and limit parameters.
func (db *DB) List(ctx context.Context, page, limit int) ([]model.Video, error) {
	collection := db.mongo.Collection(Collection)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	skip := int64((page - 1) * limit)
	lim := int64(limit)
	opts := options.Find().SetSkip(skip).SetLimit(lim).SetSort(bson.D{{Key: "createdAt", Value: 1}})
	cur, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		db.logger.WithError(err).Error("error listing videos")

		return nil, err
	}
	defer cur.Close(ctx)

	var videos []model.Video
	for cur.Next(ctx) {
		var v video
		if err := cur.Decode(&v); err != nil {
			return nil, err
		}
		videos = append(videos, v.toModel())
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return videos, nil
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
