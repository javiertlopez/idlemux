package repository

import (
	"context"

	"github.com/javiertlopez/awesome/pkg/model"
)

// AssetRepo interface
type AssetRepo interface {
	Create(ctx context.Context, source string, public bool) (string, error)
	GetByID(ctx context.Context, id string) (model.Asset, error)
}

// VideoRepo interface
type VideoRepo interface {
	Create(ctx context.Context, anyVideo model.Video) (model.Video, error)
	GetByID(ctx context.Context, id string) (model.Video, error)
}
