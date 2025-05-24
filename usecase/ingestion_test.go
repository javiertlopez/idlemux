package usecase

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/javiertlopez/awesome/errorcodes"
	"github.com/javiertlopez/awesome/model"
)

// Test_Ingestion tests the Ingestion constructor function
func Test_Ingestion(t *testing.T) {
	logger := logrus.New()
	logger.Out = io.Discard
	assets := NewMockAssets(t)
	videos := NewMockVideos(t)

	usecase := Ingestion(assets, videos, logger)

	assert.NotNil(t, usecase)
	assert.Equal(t, assets, usecase.assets)
	assert.Equal(t, videos, usecase.videos)
	assert.Equal(t, logger, usecase.logger)
}

func Test_ingestion_Create(t *testing.T) {
	logger := logrus.New()
	logger.Out = io.Discard

	assetId := "4e5bf8f2-9c50-4576-b9d4-1d1fd0705885"
	id := "bc7acb34-a7e6-4eac-87bf-8d01ad06b330"
	asset := model.Asset{
		ID: assetId,
	}
	expected := model.Video{
		ID:          id,
		Title:       "Some Might Say",
		Description: "(What's the Story) Morning Glory?",
		SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
		Asset: &model.Asset{
			ID: assetId,
		},
	}

	type args struct {
		ctx      context.Context
		anyVideo model.Video
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
			name: "Assets creation failed",
			args: args{
				ctx:      context.Background(),
				anyVideo: model.Video{},
			},
			mocks:   mockReturns{},
			want:    model.Video{},
			wantErr: true,
			err:     errorcodes.ErrVideoUnprocessable,
		},
		{
			name: "Video with public source",
			args: args{
				ctx: context.Background(),
				anyVideo: model.Video{
					Title:       "Some Might Say",
					Description: "(What's the Story) Morning Glory?",
					SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
					Policy:      "public",
				},
			},
			mocks: mockReturns{
				assetResp: asset,
				assetErr:  nil,
				videoResp: expected,
				videoErr:  nil,
			},
			want:    expected,
			wantErr: false,
			err:     nil,
		},
		{
			name: "Video with signed source",
			args: args{
				ctx: context.Background(),
				anyVideo: model.Video{
					Title:       "Some Might Say",
					Description: "(What's the Story) Morning Glory?",
					SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
					Policy:      "signed",
				},
			},
			mocks: mockReturns{
				assetResp: asset,
				assetErr:  nil,
				videoResp: expected,
				videoErr:  nil,
			},
			want:    expected,
			wantErr: false,
			err:     nil,
		},
		{
			name: "Video without policy",
			args: args{
				ctx: context.Background(),
				anyVideo: model.Video{
					Title:       "Some Might Say",
					Description: "(What's the Story) Morning Glory?",
					SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
				},
			},
			mocks: mockReturns{
				assetResp: model.Asset{},
				assetErr:  nil,
				videoResp: model.Video{},
				videoErr:  nil,
			},
			want:    model.Video{},
			wantErr: true,
			err:     errorcodes.ErrIngestionFailed,
		},
		{
			name: "Video without source URL",
			args: args{
				ctx: context.Background(),
				anyVideo: model.Video{
					Title:       "Some Might Say",
					Description: "(What's the Story) Morning Glory?",
					Policy:      "public",
				},
			},
			mocks: mockReturns{
				assetResp: model.Asset{},
				assetErr:  nil,
				videoResp: model.Video{},
				videoErr:  nil,
			},
			want:    model.Video{},
			wantErr: false,
			err:     nil,
		},
		{
			name: "Video ingestion failed",
			args: args{
				ctx: context.Background(),
				anyVideo: model.Video{
					Title:       "Some Might Say",
					Description: "(What's the Story) Morning Glory?",
					SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
					Policy:      "public",
				},
			},
			mocks: mockReturns{
				assetResp: model.Asset{},
				assetErr:  errors.New("asset creation failed"),
				videoResp: model.Video{},
				videoErr:  nil,
			},
			want:    model.Video{},
			wantErr: true,
			err:     errorcodes.ErrIngestionFailed,
		},
		{
			name: "Video creation failed",
			args: args{
				ctx: context.Background(),
				anyVideo: model.Video{
					Title:       "Some Might Say",
					Description: "(What's the Story) Morning Glory?",
					SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
					Policy:      "public",
				},
			},
			mocks: mockReturns{
				assetResp: asset,
				assetErr:  nil,
				videoResp: model.Video{},
				videoErr:  errors.New("asset creation failed"),
			},
			want:    model.Video{},
			wantErr: true,
			err:     errors.New("asset creation failed"),
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

			usecase := &ingestion{
				assets,
				videos,
				testLogger,
			}

			// Set more structured mock expectations based on the test case properties
			if len(tt.args.anyVideo.Title) == 0 || len(tt.args.anyVideo.Description) == 0 {
				// For validation error test cases, no mock expectations needed
				// The function will return early due to validation failure
			} else if len(tt.args.anyVideo.SourceURL) > 0 {
				// If there's a source URL, we expect assets.Create to be called
				if tt.args.anyVideo.Policy == "public" || tt.args.anyVideo.Policy == "signed" {
					isPublic := tt.args.anyVideo.Policy == "public"
					assets.On("Create", tt.args.ctx, tt.args.anyVideo.SourceURL, isPublic).Return(asset, tt.mocks.assetErr)

					// Only expect videos.Create if assets.Create doesn't error
					if tt.mocks.assetErr == nil {
						videos.On("Create", tt.args.ctx, mock.AnythingOfType("model.Video")).Return(tt.mocks.videoResp, tt.mocks.videoErr)
					}
				}
				// For invalid policy, no mock expectations needed as function returns early
			} else {
				// For videos without source URL, only videos.Create is called
				videos.On("Create", tt.args.ctx, mock.AnythingOfType("model.Video")).Return(tt.mocks.videoResp, tt.mocks.videoErr)
			}

			got, err := usecase.Create(tt.args.ctx, tt.args.anyVideo)

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
			}
		})

	}
}
