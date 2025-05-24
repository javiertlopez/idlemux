package usecase

import (
	"context"

	"github.com/javiertlopez/awesome/model"
)

// Assets interface
type Assets interface {
	Create(ctx context.Context, source string, public bool) (model.Asset, error)
	GetByID(ctx context.Context, id string) (model.Asset, error)
}

// Videos interface
type Videos interface {
	Create(ctx context.Context, anyVideo model.Video) (model.Video, error)
	GetByID(ctx context.Context, id string) (model.Video, error)
	List(ctx context.Context, page, limit int) ([]model.Video, error)
}
