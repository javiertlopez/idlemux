package muxinc

import (
	muxgo "github.com/muxinc/mux-go/v5"
	"github.com/sirupsen/logrus"
)

// Assets struct
type assets struct {
	logger    *logrus.Logger
	mux       *muxgo.APIClient
	keyID     string
	keySecret string
	test      bool
}

// Config struct
type Config struct {
	KeyID     string
	KeySecret string
	Test      bool
}

// New returns an asset implementation (mux.com)
func New(
	l *logrus.Logger,
	m *muxgo.APIClient,
	cfg Config,
) *assets {
	return &assets{
		logger:    l,
		mux:       m,
		keyID:     cfg.KeyID,
		keySecret: cfg.KeySecret,
		test:      cfg.Test,
	}
}
