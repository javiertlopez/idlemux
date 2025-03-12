package usecase

import (
	"context"

	guuid "github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/javiertlopez/awesome/errorcodes"
	"github.com/javiertlopez/awesome/model"
)

type delivery struct {
	assets Assets
	videos Videos
	logger *logrus.Logger
}

// Delivery returns the usecase implementation
func Delivery(
	a Assets,
	v Videos,
	l *logrus.Logger,
) delivery {
	return delivery{
		assets: a,
		videos: v,
		logger: l,
	}
}

// GetByID methods
func (u delivery) GetByID(ctx context.Context, id string) (model.Video, error) {
	// Validate UUID format
	if _, err := guuid.Parse(id); err != nil {
		return model.Video{}, errorcodes.ErrInvalidID
	}

	response, err := u.videos.GetByID(ctx, id)

	if err != nil {
		return model.Video{}, err
	}

	// If video document contains an Asset ID, retrieve the information
	if response.Asset != nil {
		asset, err := u.assets.GetByID(ctx, response.Asset.ID)
		if err != nil {
			u.logger.Error(err)
			return response, nil
		}

		response.Poster = asset.Poster
		response.Thumbnail = asset.Thumbnail
		response.Sources = asset.Sources
	}

	return response, nil
}
