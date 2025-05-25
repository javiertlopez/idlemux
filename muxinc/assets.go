package muxinc

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	muxgo "github.com/muxinc/mux-go/v5"
	"github.com/sirupsen/logrus"

	"github.com/javiertlopez/idlemux/model"
)

// asset struct
type asset struct {
	data muxgo.Asset
}

// Ingest send a source file url to mux.com
// Returns a string Asset ID
func (a *Assets) Create(ctx context.Context, source string, public bool) (model.Asset, error) {
	var policy []muxgo.PlaybackPolicy

	if public {
		policy = append(policy, muxgo.PUBLIC)
	} else {
		policy = append(policy, muxgo.SIGNED)
	}

	asset, err := a.mux.AssetsApi.CreateAsset(muxgo.CreateAssetRequest{
		Input: []muxgo.InputSettings{
			{
				Url: source,
			},
		},
		PlaybackPolicy: policy,
		Test:           a.test,
	})

	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"step":   "AssetsApi.CreateAsset",
			"func":   "Create",
			"source": source,
		}).Error(err.Error())

		return model.Asset{}, err
	}

	return model.Asset{
		ID: asset.Data.Id,
	}, nil
}

// GetByID retrieves an asset from Mux.com by Asset ID
func (a *Assets) GetByID(ctx context.Context, id string) (model.Asset, error) {
	response, err := a.mux.AssetsApi.GetAsset(id)

	if err != nil {
		return model.Asset{}, err
	}

	body := asset{
		data: response.Data,
	}

	asset := body.toModel()

	if len(body.data.PlaybackIds) > 0 {
		var source, poster, thumbnail string
		playbackID := body.data.PlaybackIds[0].Id

		switch body.data.PlaybackIds[0].Policy {
		case muxgo.PUBLIC:
			source = fmt.Sprintf(
				"https://stream.mux.com/%s.m3u8",
				playbackID,
			)

			poster = fmt.Sprintf(
				"https://image.mux.com/%s/thumbnail.png?width=%s&height=%s&smart_crop=true&time=%s",
				playbackID,
				"1920",
				"1080",
				"7",
			)

			thumbnail = fmt.Sprintf(
				"https://image.mux.com/%s/thumbnail.png?width=%s&height=%s&smart_crop=true&time=%s",
				playbackID,
				"640",
				"360",
				"7",
			)
		case muxgo.SIGNED:
			token, err := a.signURL(playbackID, "v", asset.Duration, 0, 0)
			if err != nil {
				return model.Asset{}, err
			}

			source = fmt.Sprintf(
				"https://stream.mux.com/%s.m3u8?token=%s",
				playbackID,
				token,
			)

			token, err = a.signURL(playbackID, "t", asset.Duration, 1920, 1080)
			if err != nil {
				return model.Asset{}, err
			}

			poster = fmt.Sprintf(
				"https://image.mux.com/%s/thumbnail.png?token=%s",
				playbackID,
				token,
			)

			token, err = a.signURL(playbackID, "t", asset.Duration, 640, 360)
			if err != nil {
				return model.Asset{}, err
			}

			thumbnail = fmt.Sprintf(
				"https://image.mux.com/%s/thumbnail.png?token=%s",
				playbackID,
				token,
			)
		}

		asset.Poster = poster
		asset.Thumbnail = thumbnail
		asset.Sources = []model.Source{
			{
				Source: source,
				Type:   "application/x-mpegURL",
			},
		}
	}

	return asset, nil
}

func (a *asset) toModel() model.Asset {
	return model.Asset{
		ID:                  a.data.Id,
		CreatedAt:           a.data.CreatedAt,
		Status:              a.data.Status,
		Duration:            a.data.Duration,
		MaxStoredResolution: a.data.MaxStoredResolution,
		MaxStoredFrameRate:  a.data.MaxStoredFrameRate,
		AspectRatio:         a.data.AspectRatio,
		Passthrough:         a.data.Passthrough,
	}
}

func (a *Assets) signURL(
	playbackID string,
	audience string,
	duration float64,
	width int,
	height int,
) (string, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(a.keySecret)
	if err != nil {
		return "", err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(decodedKey)
	if err != nil {
		return "", err
	}

	round := int(duration * 1.6)

	claims := jwt.MapClaims{
		"sub": playbackID,
		"aud": audience,
		"exp": time.Now().Add(time.Second * time.Duration(round)).Unix(),
		"kid": a.keyID,
	}

	if audience == "t" {
		claims["time"] = 7
		claims["width"] = width
		claims["height"] = height
		claims["smart_crop"] = true
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
