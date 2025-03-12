package usecase

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/javiertlopez/awesome/errorcodes"
	"github.com/javiertlopez/awesome/model"
	"github.com/javiertlopez/awesome/usecase/mocks"
)

// Generate mocks
// mockery --keeptree --name=Videos --dir=usecase --output=usecase/mocks

func Test_delivery_GetByID(t *testing.T) {
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
					ID:        uuid,
					Poster:    "https://image.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg/thumbnail.png?width=1920\u0026height=1080\u0026smart_crop=true\u0026time=7",
					Thumbnail: "https://image.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg/thumbnail.png?width=640\u0026height=360\u0026smart_crop=true\u0026time=7",
					Sources: []model.Source{
						{
							Source: "https://stream.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg.m3u8",
							Type:   "application/x-mpegURL",
						},
					},
				},
				videoErr: nil,
			},
			want: model.Video{
				ID:        uuid,
				Poster:    "https://image.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg/thumbnail.png?width=1920\u0026height=1080\u0026smart_crop=true\u0026time=7",
				Thumbnail: "https://image.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg/thumbnail.png?width=640\u0026height=360\u0026smart_crop=true\u0026time=7",
				Sources: []model.Source{
					{
						Source: "https://stream.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg.m3u8",
						Type:   "application/x-mpegURL",
					},
				},
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
			assets := &mocks.Assets{}
			videos := &mocks.Videos{}
			usecase := &delivery{
				assets,
				videos,
				logger,
			}

			assets.On("GetByID", tt.args.ctx, tt.args.id).Return(tt.mocks.assetResp, tt.mocks.assetErr)
			videos.On("GetByID", tt.args.ctx, tt.args.id).Return(tt.mocks.videoResp, tt.mocks.videoErr)

			got, err := usecase.GetByID(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.err, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
