package muxinc

import (
	"context"

	"github.com/javiertlopez/awesome/pkg/errorcodes"
	"github.com/javiertlopez/awesome/pkg/model"
	"github.com/stretchr/testify/mock"
)

// MockedAssets struct
type MockedAssets struct {
	mock.Mock
}

// Create mocked method
func (m *MockedAssets) Create(ctx context.Context, source string, public bool) (string, error) {
	ValidURL := "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4"
	switch source {
	case ValidURL:
		return "5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg", nil
	}
	return "", errorcodes.ErrIngestionFailed
}

// GetByID mocked method
func (m *MockedAssets) GetByID(ctx context.Context, id string) (model.Asset, error) {
	ID := "5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg"

	switch id {
	case ID:
		return model.Asset{
			ID:        ID,
			Poster:    "https://image.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg/thumbnail.png?width=1920&height=1080&smart_crop=true&time=7",
			Thumbnail: "https://image.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg/thumbnail.png?width=640&height=360&smart_crop=true&time=7",
			Sources: []model.Source{
				{
					Source: "https://stream.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg.m3u8",
					Type:   "application/x-mpegURL",
				},
			},
		}, nil
	}

	return model.Asset{}, errorcodes.ErrAssetNotFound
}
