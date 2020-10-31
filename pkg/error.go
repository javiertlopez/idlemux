package awesome

import "errors"

// ErrVideoNotFound definition
var ErrVideoNotFound = errors.New("video not found")

// ErrAssetNotFound definition
var ErrAssetNotFound = errors.New("asset not found")

// ErrIngestionFailed definition
var ErrIngestionFailed = errors.New("ingestion failed")
