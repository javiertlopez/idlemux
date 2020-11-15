package video

import (
	"context"

	awesome "github.com/javiertlopez/awesome/pkg"
	"github.com/stretchr/testify/mock"
)

type MockedVideos struct {
	mock.Mock
}

func (m *MockedVideos) Insert(ctx context.Context, anyVideo *awesome.Video) (*awesome.Video, error) {
	uuid := "fcdf5f4e-b086-4b52-8714-bf3623186185"

	response := &awesome.Video{
		ID:          &uuid,
		Title:       anyVideo.Title,
		Description: anyVideo.Description,
	}

	if anyVideo.Asset != nil {
		response.Asset = &awesome.Asset{
			ID: anyVideo.Asset.ID,
		}
	}

	return response, nil
}

func (m *MockedVideos) GetByID(ctx context.Context, id string) (*awesome.Video, error) {
	ID := "fcdf5f4e-b086-4b52-8714-bf3623186185"
	IDWithSourceFile := "a9200233-9b62-489c-9cbc-bb37f2922804"

	switch id {
	case ID:
		return &awesome.Video{
			ID:          &ID,
			Title:       "Some Might Say",
			Description: "Oasis song from (What's the Story) Morning Glory? album.",
		}, nil
	case IDWithSourceFile:
		return &awesome.Video{
			ID:          &IDWithSourceFile,
			Title:       "Some Might Say",
			Description: "Oasis song from (What's the Story) Morning Glory? album.",
			Asset: &awesome.Asset{
				ID: "5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg",
			},
		}, nil
	}

	return nil, awesome.ErrVideoNotFound
}
