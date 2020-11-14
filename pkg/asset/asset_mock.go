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
			ID: ID,
		}, nil
	}

	return nil, awesome.ErrAssetNotFound
}
