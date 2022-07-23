package usecase

import (
	"context"

	"github.com/javiertlopez/awesome/model"
)

type delivery struct {
	assets Assets
	videos Videos
}

// Delivery returns the usecase implementation
func Delivery(
	a Assets,
	v Videos,
) delivery {
	return delivery{
		assets: a,
		videos: v,
	}
}

// GetByID methods
func (u delivery) GetByID(ctx context.Context, id string) (model.Video, error) {
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
