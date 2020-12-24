package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	mocks "github.com/javiertlopez/awesome/mocks/repository"
	"github.com/javiertlopez/awesome/model"
	"github.com/stretchr/testify/mock"
)

func Test_ingestion_Create(t *testing.T) {
	uuid := "4e5bf8f2-9c50-4576-b9d4-1d1fd0705885"
	expected := model.Video{
		Title:       "Some Might Say",
		Description: "(What's the Story) Morning Glory?",
		SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
		Asset: &model.Asset{
			ID: uuid,
		},
	}
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
			"Video with public source",
			args{
				ctx: context.Background(),
				anyVideo: model.Video{
					Title:       "Some Might Say",
					Description: "(What's the Story) Morning Glory?",
					SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
					Policy:      "public",
				},
			},
			expected,
			false,
		},
		{
			"Video with signed source",
			args{
				ctx: context.Background(),
				anyVideo: model.Video{
					Title:       "Some Might Say",
					Description: "(What's the Story) Morning Glory?",
					SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
					Policy:      "signed",
				},
			},
			expected,
			false,
		},
		{
			"Video without policy",
			args{
				ctx: context.Background(),
				anyVideo: model.Video{
					Title:       "Some Might Say",
					Description: "(What's the Story) Morning Glory?",
					SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
				},
			},
			model.Video{},
			true,
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
			u := &videos{
				ar,
				vr,
			}

			if tt.wantErr {
				ar.On("Create", tt.args.ctx, tt.args.anyVideo.SourceURL, mock.Anything).Return(uuid, errors.New("failed"))
				vr.On("Create", tt.args.ctx, tt.args.anyVideo).Return(model.Video{}, errors.New("failed"))
			} else {
				ar.On("Create", tt.args.ctx, tt.args.anyVideo.SourceURL, mock.Anything).Return(uuid, nil)
				vr.On("Create", tt.args.ctx, mock.Anything).Return(tt.want, nil)
			}
			got, err := u.Create(tt.args.ctx, tt.args.anyVideo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ingestion.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ingestion.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
