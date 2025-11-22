package usecase

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/javiertlopez/idlemux/errorcodes"
	"github.com/javiertlopez/idlemux/model"
)

// TestDelivery tests the Delivery constructor function
func TestDelivery(t *testing.T) {
	logger := logrus.New()
	logger.Out = io.Discard
	assets := NewMockAssets(t)
	videos := NewMockVideos(t)

	usecase := Delivery(assets, videos, logger)

	assert.NotNil(t, usecase)
	assert.Equal(t, assets, usecase.assets)
	assert.Equal(t, videos, usecase.videos)
	assert.Equal(t, logger, usecase.logger)
}

func TestDelivery_GetByID(t *testing.T) {
	logger := logrus.New()
	logger.Out = io.Discard

	uuid := "4e5bf8f2-9c50-4576-b9d4-1d1fd0705885"
	asset := model.Asset{
		ID:        uuid,
		Poster:    "https://image.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg/thumbnail.png?width=1920\u0026height=1080\u0026smart_crop=true\u0026time=7",
		Thumbnail: "https://image.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg/thumbnail.png?width=640\u0026height=360\u0026smart_crop=true\u0026time=7",
		Sources: []model.Source{
			{
				Source: "https://stream.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg.m3u8",
				Type:   "application/x-mpegURL",
			},
		},
	}

	type args struct {
		ctx context.Context
		id  string
	}

	type mockReturns struct {
		assetResp model.Asset
		assetErr  error
		videoResp model.Video
		videoErr  error
	}

	tests := []struct {
		name    string
		args    args
		mocks   mockReturns
		want    model.Video
		wantErr bool
		err     error
	}{
		{
			name: "Only video",
			args: args{
				ctx: context.Background(),
				id:  uuid,
			},
			mocks: mockReturns{
				assetResp: asset,
				assetErr:  nil,
				videoResp: model.Video{},
				videoErr:  nil,
			},
			want:    model.Video{},
			wantErr: false,
			err:     nil,
		},
		{
			name: "Asset service error",
			args: args{
				ctx: context.Background(),
				id:  uuid,
			},
			mocks: mockReturns{
				assetResp: model.Asset{},
				assetErr:  errors.New("asset service failed"),
				videoResp: model.Video{
					Asset: &model.Asset{
						ID: uuid,
					},
				},
				videoErr: nil,
			},
			want: model.Video{
				Asset: &model.Asset{
					ID: uuid,
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "Error",
			args: args{
				ctx: context.Background(),
				id:  uuid,
			},
			mocks: mockReturns{
				assetResp: asset,
				assetErr:  nil,
				videoResp: model.Video{},
				videoErr:  errors.New("failed"),
			},
			want:    model.Video{},
			wantErr: true,
			err:     errors.New("failed"),
		},
		{
			name: "With asset",
			args: args{
				ctx: context.Background(),
				id:  uuid,
			},
			mocks: mockReturns{
				assetResp: asset,
				assetErr:  nil,
				videoResp: model.Video{
					ID: uuid,
					Asset: &model.Asset{
						ID: uuid,
					},
				},
				videoErr: nil,
			},
			want: model.Video{
				ID: uuid,
				Asset: &model.Asset{
					ID: uuid,
				},
				Poster:    asset.Poster,
				Thumbnail: asset.Thumbnail,
				Sources:   asset.Sources,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "Video without asset",
			args: args{
				ctx: context.Background(),
				id:  uuid,
			},
			mocks: mockReturns{
				assetResp: asset,
				assetErr:  nil,
				videoResp: model.Video{
					ID: uuid,
					// No Asset field (nil)
				},
				videoErr: nil,
			},
			want: model.Video{
				ID: uuid,
				// Asset field should remain nil
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "With invalid ID",
			args: args{
				ctx: context.Background(),
				id:  "invalid",
			},
			mocks: mockReturns{
				assetResp: asset,
				assetErr:  nil,
				videoResp: model.Video{},
				videoErr:  nil,
			},
			want:    model.Video{},
			wantErr: true,
			err:     errorcodes.ErrInvalidID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks with strict mode to catch any unexpected calls
			assets := NewMockAssets(t)
			videos := NewMockVideos(t)

			// Create a real logger that discards output but can be inspected for coverage
			testLogger := logrus.New()
			testLogger.Out = io.Discard

			usecase := &delivery{
				assets,
				videos,
				testLogger,
			}

			// Only set expectations if the test case is expected to call the mocks
			if tt.name == "With invalid ID" {
				// No expectations for invalid ID case
			} else {
				// All other test cases call videos.GetByID
				videos.On("GetByID", tt.args.ctx, tt.args.id).Return(tt.mocks.videoResp, tt.mocks.videoErr)

				// Set up asset expectations for cases with assets
				if tt.mocks.videoResp.Asset != nil {
					// Only set up GetByID expectation if the Asset ID is not empty
					if tt.mocks.videoResp.Asset.ID != "" {
						assets.On("GetByID", tt.args.ctx, tt.mocks.videoResp.Asset.ID).Return(tt.mocks.assetResp, tt.mocks.assetErr)
					}
					// Note: If the Asset ID is empty, the code shouldn't call GetByID
				}
			}
			// Don't use AssertNotCalled here - it will be checked before the function is executed

			got, err := usecase.GetByID(tt.args.ctx, tt.args.id)

			// For invalid ID cases, verify we don't call the mocks at all
			if tt.name == "With invalid ID" {
				videos.AssertNotCalled(t, "GetByID")
				assets.AssertNotCalled(t, "GetByID")
			}

			if tt.wantErr {
				assert.Error(t, err, "Expected an error but got none")
				if tt.err != nil {
					assert.Equal(t, tt.err.Error(), err.Error(), "Error message doesn't match")
				}
				return
			}

			assert.NoError(t, err, "Got unexpected error")
			assert.Equal(t, tt.want, got, "Response doesn't match expected")

			// Additional specific assertions for certain fields
			if tt.want.Asset != nil {
				assert.Equal(t, tt.want.Asset.ID, got.Asset.ID, "Asset ID doesn't match")

				// If we expect populated assets, verify those fields too
				if tt.name == "With asset" {
					assert.Equal(t, tt.want.Poster, got.Poster, "Poster field doesn't match")
					assert.Equal(t, tt.want.Thumbnail, got.Thumbnail, "Thumbnail field doesn't match")
					assert.Equal(t, tt.want.Sources, got.Sources, "Sources field doesn't match")
				}
			}
		})
	}
}

func TestDelivery_List(t *testing.T) {
	logger := logrus.New()
	logger.Out = io.Discard

	videosList := []model.Video{
		{ID: "id1", Title: "Video 1"},
		{ID: "id2", Title: "Video 2"},
	}

	type args struct {
		ctx   context.Context
		page  int
		limit int
	}

	type mockReturns struct {
		videosResp []model.Video
		videosErr  error
	}

	tests := []struct {
		name    string
		args    args
		mocks   mockReturns
		want    []model.Video
		wantErr bool
		err     error
	}{
		{
			name: "Success",
			args: args{
				ctx:   context.Background(),
				page:  1,
				limit: 10,
			},
			mocks: mockReturns{
				videosResp: videosList,
				videosErr:  nil,
			},
			want:    videosList,
			wantErr: false,
			err:     nil,
		},
		{
			name: "Error",
			args: args{
				ctx:   context.Background(),
				page:  1,
				limit: 10,
			},
			mocks: mockReturns{
				videosResp: nil,
				videosErr:  errors.New("db error"),
			},
			want:    nil,
			wantErr: true,
			err:     errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			videos := NewMockVideos(t)
			assets := NewMockAssets(t)
			usecase := &delivery{
				assets,
				videos,
				logger,
			}

			videos.On("List", tt.args.ctx, tt.args.page, tt.args.limit).Return(tt.mocks.videosResp, tt.mocks.videosErr)

			got, err := usecase.List(tt.args.ctx, tt.args.page, tt.args.limit)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.err.Error(), err.Error()) // Compare error messages for more robust comparison
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
