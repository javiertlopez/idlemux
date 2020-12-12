package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	mocks "github.com/javiertlopez/awesome/pkg/mocks/repository"
	"github.com/javiertlopez/awesome/pkg/model"

	"github.com/stretchr/testify/mock"
)

func Test_videos_Create(t *testing.T) {
	uuid := "4e5bf8f2-9c50-4576-b9d4-1d1fd0705885"
	type args struct {
		ctx      context.Context
		anyVideo model.Video
	}
	tests := []struct {
		name    string
		args    args
		want    model.Video
		wantErr bool
	}{
		{
			"Video Unprocessable",
			args{
				ctx:      context.Background(),
				anyVideo: model.Video{},
			},
			model.Video{},
			true,
		},
		{
			"Video with source",
			args{
				ctx: context.Background(),
				anyVideo: model.Video{
					Title:       "Some Might Say",
					Description: "(What's the Story) Morning Glory?",
					SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
				},
			},
			model.Video{
				Title:       "Some Might Say",
				Description: "(What's the Story) Morning Glory?",
				SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
				Asset: &model.Asset{
					ID: uuid,
				},
			},
			false,
		},
		{
			"Video with source (error)",
			args{
				ctx: context.Background(),
				anyVideo: model.Video{
					Title:       "Some Might Say",
					Description: "(What's the Story) Morning Glory?",
					SourceURL:   "fakeurl",
				},
			},
			model.Video{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := &mocks.AssetRepo{}
			vr := &mocks.VideoRepo{}
			v := &videos{
				ar,
				vr,
			}

			if tt.wantErr {
				ar.On("Create", tt.args.ctx, tt.args.anyVideo.SourceURL, true).Return(uuid, errors.New("failed"))
				vr.On("Create", tt.args.ctx, tt.args.anyVideo).Return(model.Video{}, errors.New("failed"))
			} else {
				ar.On("Create", tt.args.ctx, tt.args.anyVideo.SourceURL, true).Return(uuid, nil)
				vr.On("Create", tt.args.ctx, mock.Anything).Return(tt.want, nil)
			}

			got, err := v.Create(tt.args.ctx, tt.args.anyVideo)
			if (err != nil) != tt.wantErr {
				t.Errorf("videos.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("videos.Create() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_videos_GetByID(t *testing.T) {
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
				ID: &uuid,
				Asset: &model.Asset{
					ID: "4e5bf8f2-9c50-4576-b9d4-1d1fd0705885",
				},
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
			ar := &mocks.AssetRepo{}
			vr := &mocks.VideoRepo{}
			v := &videos{
				ar,
				vr,
			}

			if tt.wantErr {
				ar.On("GetByID", tt.args.ctx, tt.args.id).Return(model.Asset{}, errors.New("failed"))
				vr.On("GetByID", tt.args.ctx, tt.args.id).Return(model.Video{}, errors.New("failed"))
			} else {
				ar.On("GetByID", tt.args.ctx, tt.args.id).Return(asset, nil)
				vr.On("GetByID", tt.args.ctx, tt.args.id).Return(tt.want, nil)
			}

			got, err := v.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("videos.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("videos.GetByID() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
