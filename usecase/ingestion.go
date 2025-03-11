package usecase

import (
	"context"

	"github.com/javiertlopez/awesome/errorcodes"
	"github.com/javiertlopez/awesome/model"
)

type ingestion struct {
	assets Assets
	videos Videos
}

// Ingestion returns the usecase implementation
func Ingestion(
	a Assets,
	v Videos,
) ingestion {
	return ingestion{
		assets: a,
		videos: v,
	}
}

// Create method
func (u ingestion) Create(ctx context.Context, anyVideo model.Video) (model.Video, error) {
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

		asset, err := u.assets.Create(ctx, anyVideo.SourceURL, isPublic)
		if err == nil {
			anyVideo.Asset = &asset
		}
	}

	response, err := u.videos.Create(ctx, anyVideo)

	if err != nil {
		return model.Video{}, err
	}

	return response, nil
}
