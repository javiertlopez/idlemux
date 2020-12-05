package axiom

import (
	"context"

	"github.com/javiertlopez/awesome/pkg/errorcodes"
	"github.com/javiertlopez/awesome/pkg/model"

	"github.com/stretchr/testify/mock"
)

// MockedVideos struct
type MockedVideos struct {
	mock.Mock
}

// Create mocked method
func (m *MockedVideos) Create(ctx context.Context, anyVideo model.Video) (model.Video, error) {
	uuid := "fcdf5f4e-b086-4b52-8714-bf3623186185"

	response := model.Video{
		ID:          &uuid,
		Title:       anyVideo.Title,
		Description: anyVideo.Description,
	}

	if anyVideo.Asset != nil {
		response.Asset = &model.Asset{
			ID: anyVideo.Asset.ID,
		}
	}

	return response, nil
}

// GetByID mocked method
func (m *MockedVideos) GetByID(ctx context.Context, id string) (*model.Video, error) {
	ID := "fcdf5f4e-b086-4b52-8714-bf3623186185"
	IDWithSourceFile := "a9200233-9b62-489c-9cbc-bb37f2922804"

	switch id {
	case ID:
		return &model.Video{
			ID:          &ID,
			Title:       "Some Might Say",
			Description: "Oasis song from (What's the Story) Morning Glory? album.",
		}, nil
	case IDWithSourceFile:
		return &model.Video{
			ID:          &IDWithSourceFile,
			Title:       "Some Might Say",
			Description: "Oasis song from (What's the Story) Morning Glory? album.",
			Asset: &model.Asset{
				ID: "5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg",
			},
		}, nil
	}

	return nil, errorcodes.ErrVideoNotFound
}
