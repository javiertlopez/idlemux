package controller

import (
	"context"

	"github.com/javiertlopez/awesome/model"
)

// Delivery usecase
type Delivery interface {
	GetByID(ctx context.Context, id string) (model.Video, error)
	List(ctx context.Context, page, limit int) ([]model.Video, error)
}

// Ingestion usecase
type Ingestion interface {
	Create(ctx context.Context, anyVideo model.Video) (model.Video, error)
}

// controller struct holds the usecase
type controller struct {
	commit    string
	version   string
	delivery  Delivery
	ingestion Ingestion
}

// New returns a controller
func New(
	commit string,
	version string,
	delivery Delivery,
	ingestion Ingestion,
) controller {
	return controller{
		commit: commit,

		version:   version,
		delivery:  delivery,
		ingestion: ingestion,
	}
}
