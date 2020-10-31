package main

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockedAssets struct {
	mock.Mock
}

func (m *MockedAssets) Ingest(ctx context.Context, source string) (string, error) {
	ValidURL := "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4"
	switch source {
	case ValidURL:
		return "5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg", nil
	}
	return "", ErrIngestionFailed
}

func (m *MockedAssets) GetByID(ctx context.Context, id string) (*Asset, error) {
	ID := "5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg"

	switch id {
	case ID:
		return &Asset{
			ID: ID,
		}, nil
	}

	return nil, ErrAssetNotFound
}
