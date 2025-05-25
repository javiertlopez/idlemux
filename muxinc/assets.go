package muxinc

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	muxgo "github.com/muxinc/mux-go/v5"

	"github.com/javiertlopez/idlemux/model"
)

// Constants for image dimensions
const (
	posterWidth     = 1920
	posterHeight    = 1080
	thumbnailWidth  = 640
	thumbnailHeight = 360
	imageTime       = "7"
)

// asset struct
type asset struct {
	data muxgo.Asset
}

// Ingest send a source file url to mux.com
// Returns a string Asset ID
func (a *assets) Create(ctx context.Context, source string, public bool) (model.Asset, error) {
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
		a.logger.WithError(err).Error("error creating asset")

		return model.Asset{}, err
	}

	return model.Asset{
		ID: asset.Data.Id,
	}, nil
}

// GetByID retrieves an asset from Mux.com by Asset ID
func (a *assets) GetByID(ctx context.Context, id string) (model.Asset, error) {
	response, err := a.mux.AssetsApi.GetAsset(id)
	if err != nil {
		a.logger.WithError(err).Error("error retrieving asset by ID")

		return model.Asset{}, err
	}

	body := asset{
		data: response.Data,
	}

	asset := body.toModel()

	if len(body.data.PlaybackIds) > 0 {
		playbackID := body.data.PlaybackIds[0].Id
		policy := body.data.PlaybackIds[0].Policy

		if err := a.hydrateAssetURLs(playbackID, policy, asset.Duration, &asset); err != nil {
			a.logger.WithError(err).Error("error generating asset URLs")

			return model.Asset{}, err
		}
	}

	return asset, nil
}

// hydrateAssetURLs adds source, poster, and thumbnail URLs to the asset
func (a *assets) hydrateAssetURLs(playbackID string, policy muxgo.PlaybackPolicy, duration float64, asset *model.Asset) error {
	var source, poster, thumbnail string
	var err error

	switch policy {
	case muxgo.PUBLIC:
		source, poster, thumbnail = a.generatePublicURLs(playbackID)
	case muxgo.SIGNED:
		source, poster, thumbnail, err = a.generateSignedURLs(playbackID, duration)
		if err != nil {
			return err
		}
	}

	asset.Poster = poster
	asset.Thumbnail = thumbnail
	asset.Sources = []model.Source{
		{
			Source: source,
			Type:   "application/x-mpegURL",
		},
	}

	return nil
}

// generatePublicURLs creates public URLs for the asset
func (a *assets) generatePublicURLs(playbackID string) (source, poster, thumbnail string) {
	source = fmt.Sprintf(
		"https://stream.mux.com/%s.m3u8",
		playbackID,
	)

	poster = fmt.Sprintf(
		"https://image.mux.com/%s/thumbnail.png?width=%d&height=%d&smart_crop=true&time=%s",
		playbackID,
		posterWidth,
		posterHeight,
		imageTime,
	)

	thumbnail = fmt.Sprintf(
		"https://image.mux.com/%s/thumbnail.png?width=%d&height=%d&smart_crop=true&time=%s",
		playbackID,
		thumbnailWidth,
		thumbnailHeight,
		imageTime,
	)

	return source, poster, thumbnail
}

// generateSignedURLs creates signed URLs for the asset
func (a *assets) generateSignedURLs(playbackID string, duration float64) (source, poster, thumbnail string, err error) {
	// Generate video token
	videoToken, err := a.signURL(playbackID, "v", duration, 0, 0)
	if err != nil {
		return "", "", "", fmt.Errorf("error signing URL for video playback: %w", err)
	}

	// Generate poster token
	posterToken, err := a.signURL(playbackID, "t", duration, posterWidth, posterHeight)
	if err != nil {
		return "", "", "", fmt.Errorf("error signing URL for poster: %w", err)
	}

	// Generate thumbnail token
	thumbnailToken, err := a.signURL(playbackID, "t", duration, thumbnailWidth, thumbnailHeight)
	if err != nil {
		return "", "", "", fmt.Errorf("error signing URL for thumbnail: %w", err)
	}

	source = fmt.Sprintf(
		"https://stream.mux.com/%s.m3u8?token=%s",
		playbackID,
		videoToken,
	)

	poster = fmt.Sprintf(
		"https://image.mux.com/%s/thumbnail.png?token=%s",
		playbackID,
		posterToken,
	)

	thumbnail = fmt.Sprintf(
		"https://image.mux.com/%s/thumbnail.png?token=%s",
		playbackID,
		thumbnailToken,
	)

	return source, poster, thumbnail, nil
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

func (a *assets) signURL(
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
