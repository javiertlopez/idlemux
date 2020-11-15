package asset

import (
	"context"

	awesome "github.com/javiertlopez/awesome/pkg"
	"github.com/stretchr/testify/mock"
)

// MockedAssets struct
type MockedAssets struct {
	mock.Mock
}

// Ingest mocked method
func (m *MockedAssets) Ingest(ctx context.Context, source string, public bool) (string, error) {
	ValidURL := "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4"
	switch source {
	case ValidURL:
		return "5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg", nil
	}
	return "", awesome.ErrIngestionFailed
}

// GetByID mocked method
func (m *MockedAssets) GetByID(ctx context.Context, id string) (*awesome.Asset, error) {
	ID := "5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg"

	switch id {
	case ID:
		return &awesome.Asset{
			ID:        ID,
			Poster:    "https://image.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg/thumbnail.png?width=1920&height=1080&smart_crop=true&time=7",
			Thumbnail: "https://image.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg/thumbnail.png?width=640&height=360&smart_crop=true&time=7",
			Sources: []awesome.Source{
				{
					ID:     "5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg",
					Policy: "public",
					Source: "https://stream.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg.m3u8",
					Type:   "application/x-mpegURL",
				},
			},
		}, nil
	}

	return nil, awesome.ErrAssetNotFound
}
