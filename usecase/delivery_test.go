package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/javiertlopez/awesome/model"
	"github.com/javiertlopez/awesome/usecase/mocks"
)

// Generate mocks
// mockery --keeptree --name=Videos --dir=usecase --output=usecase/mocks

func Test_delivery_GetByID(t *testing.T) {
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
	tests := []struct {
		name    string
		args    args
		want    model.Video
		wantErr bool
	}{
		{
			"Only video",
			args{
				ctx: context.Background(),
				id:  uuid,
			},
			model.Video{},
			false,
		},
		{
			"Error",
			args{
				ctx: context.Background(),
				id:  uuid,
			},
			model.Video{},
			true,
		},
		{
			"With asset",
			args{
				ctx: context.Background(),
				id:  uuid,
			},
			model.Video{
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
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assets := &mocks.Assets{}
			videos := &mocks.Videos{}
			usecase := &delivery{
				assets,
				videos,
			}

			if tt.wantErr {
				assets.On("GetByID", tt.args.ctx, tt.args.id).Return(model.Asset{}, errors.New("failed"))
				videos.On("GetByID", tt.args.ctx, tt.args.id).Return(model.Video{}, errors.New("failed"))
			} else {
				assets.On("GetByID", tt.args.ctx, tt.args.id).Return(asset, nil)
				videos.On("GetByID", tt.args.ctx, tt.args.id).Return(tt.want, nil)
			}

			got, err := usecase.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("delivery.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("delivery.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
