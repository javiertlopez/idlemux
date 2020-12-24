package usecase

import (
	"context"

	"github.com/javiertlopez/awesome/errorcodes"
	"github.com/javiertlopez/awesome/model"
	"github.com/javiertlopez/awesome/repository"
)

// Videos usecase
type Videos interface {
	Create(ctx context.Context, anyVideo model.Video) (model.Video, error)
	GetByID(ctx context.Context, id string) (model.Video, error)
}

type videos struct {
	assets repository.AssetRepo
	videos repository.VideoRepo
}

// NewVideoUseCase returns the usecase implementation
func NewVideoUseCase(
	a repository.AssetRepo,
	v repository.VideoRepo,
) Videos {
	return &videos{
		assets: a,
		videos: v,
	}
}

// Create method
func (v *videos) Create(ctx context.Context, anyVideo model.Video) (model.Video, error) {
	// Title and Description are mandatory fields
	if len(anyVideo.Title) == 0 || len(anyVideo.Description) == 0 {
		return model.Video{}, errorcodes.ErrVideoUnprocessable
	}

	// If body contains a Source File URL, send it to Ingestion
	if len(anyVideo.SourceURL) > 0 {
		var isPublic bool
		switch anyVideo.Policy {
		case "public":
			isPublic = true
		case "signed":
			isPublic = false
		default:
			return model.Video{}, errorcodes.ErrIngestionFailed
		}

		assetID, err := v.assets.Create(ctx, anyVideo.SourceURL, isPublic)
		if err == nil {
			anyVideo.Asset = &model.Asset{
				ID: assetID,
			}
		}
	}

	response, err := v.videos.Create(ctx, anyVideo)

	if err != nil {
		return model.Video{}, err
	}

	return response, nil
}

// GetByID methods
func (v *videos) GetByID(ctx context.Context, id string) (model.Video, error) {
	response, err := v.videos.GetByID(ctx, id)

	if err != nil {
		return model.Video{}, err
	}

	// If video document contains an Asset ID, retrieve the information
	if response.Asset != nil {
		asset, err := v.assets.GetByID(ctx, response.Asset.ID)
		if err == nil {
			response.Asset = &model.Asset{
				ID: asset.ID,
			}
			response.Poster = asset.Poster
			response.Thumbnail = asset.Thumbnail
			response.Sources = asset.Sources
		}
	}

	return response, nil
}
