package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/javiertlopez/idlemux/errorcodes"
	"github.com/javiertlopez/idlemux/model"
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
	if _, err := uuid.Parse(id); err != nil {
		return model.Video{}, errorcodes.ErrInvalidID
	}

	response, err := u.videos.GetByID(ctx, id)

	if err != nil {
		u.logger.WithError(err).Error(err.Error())
		return model.Video{}, err
	}

	// If video document contains an Asset ID, retrieve the information
	if response.Asset != nil {
		asset, err := u.assets.GetByID(ctx, response.Asset.ID)
		if err != nil {
			u.logger.WithError(err).Error(err.Error())
			return response, nil
		}

		response.Poster = asset.Poster
		response.Thumbnail = asset.Thumbnail
		response.Sources = asset.Sources
	}

	return response, nil
}

// List method
func (u delivery) List(ctx context.Context, page, limit int) ([]model.Video, error) {
	videos, err := u.videos.List(ctx, page, limit)
	if err != nil {
		u.logger.WithError(err).Error(err.Error())
		return nil, err
	}
	return videos, nil
}
