# idlemux

[![Go Reference](https://pkg.go.dev/badge/github.com/javiertlopez/idlemux.svg)](https://pkg.go.dev/github.com/javiertlopez/idlemux)
[![Go Report Card](https://goreportcard.com/badge/github.com/javiertlopez/idlemux)](https://goreportcard.com/report/github.com/javiertlopez/idlemux)
[![Go](https://github.com/javiertlopez/idlemux/workflows/Go/badge.svg)](https://github.com/javiertlopez/idlemux/actions)
[![codecov](https://codecov.io/gh/javiertlopez/idlemux/branch/main/graph/badge.svg?token=I8D2Z4TZX4)](https://codecov.io/gh/javiertlopez/idlemux)
[![Go Version](https://img.shields.io/github/go-mod/go-version/javiertlopez/idlemux)](https://github.com/javiertlopez/idlemux/blob/main/go.mod)
[![MIT License](https://img.shields.io/github/license/javiertlopez/idlemux)](https://github.com/javiertlopez/idlemux/blob/main/LICENSE)

## Overview

IdleMux is a REST API that stores a Video Library in MongoDB Atlas and uses Mux.com for video processing and delivery. It provides endpoints for creating, retrieving, and listing videos with a clean and consistent API design.

## Features

- Video management (create, get by ID, list with pagination)
- Video ingestion through Mux.com
- Health check and application status endpoints
- Well-documented API with OpenAPI 3.0.1 specification
- Comprehensive test suite with high code coverage

## Getting started

### Prerequisites

- Go 1.23 or 1.24 (we test against the latest two Go versions)
- MongoDB Atlas account
- Mux.com account with API credentials

### Environment Variables

Set up the following environment variables:

```
ADDR=:8080                      # Server address (default :8080)
MONGO_STRING=mongodb://...      # MongoDB connection string
MUX_TOKEN_ID=                   # Mux API token ID
MUX_TOKEN_SECRET=               # Mux API token secret
MUX_KEY_ID=                     # Mux signing key ID
MUX_KEY_SECRET=                 # Mux signing key secret
```

## Test and build

Run the test suite:

```bash
make test
```

Build the application:

```bash
make build
```

## API Documentation

The API is documented using OpenAPI 3.0.1. You can find the specification in the [openapi.yaml](./openapi.yaml) file.

### Endpoints

| Method | Path          | Description                                   |
|--------|---------------|-----------------------------------------------|
| GET    | /app/healthz  | Health check endpoint                         |
| GET    | /app/statusz  | Get application version and commit information|
| GET    | /videos       | List videos with pagination                   |
| POST   | /videos       | Create a new video                            |
| GET    | /videos/{id}  | Get a video by ID                             |

## Usage

Install as a dependency:

```bash
go get github.com/javiertlopez/idlemux
```

### Example

```go
package main

import (
	"net/http"
	"os"
	"time"

	"github.com/javiertlopez/idlemux"

	"github.com/sirupsen/logrus"
)

const (
	writeTimeout = 15 * time.Second
	readTimeout  = 15 * time.Second
	idleTimeout  = 60 * time.Second
)

var (
	application idlemux.App
	commit      string // Set during build with -ldflags
	version     string // Set during build with -ldflags
)

func main() {
	// Environment variables
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	mongoString := os.Getenv("MONGO_STRING")
	muxTokenID := os.Getenv("MUX_TOKEN_ID")
	muxTokenSecret := os.Getenv("MUX_TOKEN_SECRET")
	muxKeyID := os.Getenv("MUX_KEY_ID")
	muxKeySecret := os.Getenv("MUX_KEY_SECRET")

	// Create a logrus logger and set up the output format as JSON
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	application = idlemux.New(
		idlemux.AppConfig{
			Commit:         commit,
			Version:        version,
			MongoURI:       mongoString,
			MuxTokenID:     muxTokenID,
			MuxTokenSecret: muxTokenSecret,
			MuxKeyID:       muxKeyID,
			MuxKeySecret:   muxKeySecret,
		},
		logger,
	)

	// Create a Gorilla Mux router
	router := application.Router()

	// Create a Server instance with the router
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      router,
	}

	logger.WithFields(logrus.Fields{
		"addr":    addr,
		"version": version,
		"commit":  commit,
	}).Info("Starting server")

	// Start the server
	logger.Fatal(srv.ListenAndServe())
}
```

## Development

### Go Version Compatibility

This project is tested against the latest two Go versions (currently 1.23 and 1.24). We use GitHub Actions to ensure compatibility with these versions.

### Building with Version Information

To inject version and commit information during build:

```bash
go build -ldflags="-X 'main.version=v0.7.0' -X 'main.commit=$(git rev-parse --short HEAD)'" -o idlemux
```

### Adding API Endpoints

1. Define the endpoint in `openapi.yaml`
2. Add the handler method to the Controller interface in `router/router.go`
3. Implement the handler in the controller package
4. Register the route in `router/router.go`
5. Update tests accordingly

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

Licensed under [MIT License](LICENSE). © 2025 Hiram Torres Lopez.