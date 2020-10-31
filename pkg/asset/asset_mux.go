package asset

import (
	"context"
	"fmt"

	awesome "github.com/javiertlopez/awesome/pkg"
	muxgo "github.com/muxinc/mux-go"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// videos struct holds the logger and MongoDB client
type assets struct {
	logger *logrus.Logger
	mux    *muxgo.APIClient
}

type asset struct {
	data muxgo.Asset
}

// NewAssetService creates new an Videos service object
func NewAssetService(
	l *logrus.Logger,
	m *muxgo.APIClient,
) awesome.Assets {
	return &assets{
		logger: l,
		mux:    m,
	}
}

func (a *assets) Ingest(ctx context.Context, source string) (string, error) {
	asset, err := a.mux.AssetsApi.CreateAsset(muxgo.CreateAssetRequest{
		Input: []muxgo.InputSettings{
			muxgo.InputSettings{
				Url: source,
			},
		},
		PlaybackPolicy: []muxgo.PlaybackPolicy{muxgo.PUBLIC},
	})

	if err != nil {
		a.logger.WithFields(log.Fields{
			"step":   "AssetsApi.CreateAsset",
			"func":   "func (a *assets) Ingest",
			"source": source,
		}).Error(err.Error())

		return "", err
	}

	return asset.Data.Id, nil
}

func (a *assets) GetByID(ctx context.Context, id string) (*awesome.Asset, error) {
	response, err := a.mux.AssetsApi.GetAsset(id)

	if err != nil {
		return nil, err
	}

	body := asset{
		data: response.Data,
	}

	return body.toModel(), nil
}

func (a *asset) toModel() *awesome.Asset {
	response := &awesome.Asset{
		ID:                  a.data.Id,
		CreatedAt:           a.data.CreatedAt,
		Status:              a.data.Status,
		Duration:            a.data.Duration,
		MaxStoredResolution: a.data.MaxStoredResolution,
		MaxStoredFrameRate:  a.data.MaxStoredFrameRate,
		AspectRatio:         a.data.AspectRatio,
		Passthrough:         a.data.Passthrough,
	}

	if len(a.data.PlaybackIds) > 0 {
		playbackID := a.data.PlaybackIds[0].Id

		response.Poster = fmt.Sprintf(
			"https://image.mux.com/%s/thumbnail.png?width=%s&height=%s&smart_crop=true&time=%s",
			playbackID,
			"1920",
			"1080",
			"7",
		)

		response.Thumbnail = fmt.Sprintf(
			"https://image.mux.com/%s/thumbnail.png?width=%s&height=%s&smart_crop=true&time=%s",
			playbackID,
			"640",
			"360",
			"7",
		)

		response.Sources = []awesome.Source{
			awesome.Source{
				ID: playbackID,
				Source: fmt.Sprintf(
					"https://stream.mux.com/%s.m3u8",
					playbackID,
				),
				Type: "application/x-mpegURL",
			},
		}
	}

	return response
}
