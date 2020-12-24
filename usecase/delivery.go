package usecase

import (
	"context"

	"github.com/javiertlopez/awesome/model"
	"github.com/javiertlopez/awesome/repository"
)

// Delivery usecase
type Delivery interface {
	GetByID(ctx context.Context, id string) (model.Video, error)
}

type delivery struct {
	assets repository.AssetRepo
	videos repository.VideoRepo
}

// NewDelivery returns the usecase implementation
func NewDelivery(
	a repository.AssetRepo,
	v repository.VideoRepo,
) Delivery {
	return &delivery{
		assets: a,
		videos: v,
	}
}

// GetByID methods
func (u *delivery) GetByID(ctx context.Context, id string) (model.Video, error) {
	response, err := u.videos.GetByID(ctx, id)

	if err != nil {
		return model.Video{}, err
	}

	// If video document contains an Asset ID, retrieve the information
	if response.Asset != nil {
		asset, err := u.assets.GetByID(ctx, response.Asset.ID)
		if err == nil {
			response.Asset = nil
			response.Poster = asset.Poster
			response.Thumbnail = asset.Thumbnail
			response.Sources = asset.Sources
		}
	}

	return response, nil
}
