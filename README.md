# awesome   ![Go](https://github.com/javiertlopez/awesome/workflows/Go/badge.svg)   [![codecov](https://codecov.io/gh/javiertlopez/awesome/branch/main/graph/badge.svg?token=I8D2Z4TZX4)](https://codecov.io/gh/javiertlopez/awesome)

## Overview

Awesome is a REST API that stores a Video Library in Mongo Atlas and Mux.com for video ingestion.

## Getting started

### Requirements

- github.com/golang/go go1.18.4
- github.com/golang-jwt/jwt/v4 v4.4.2
- github.com/google/uuid v1.3.0
- github.com/gorilla/mux v1.8.0
- github.com/muxinc/mux-go/v4 v4.2.3
- github.com/sirupsen/logrus v1.9.0
- github.com/stretchr/testify v1.8.0
- go.mongodb.org/mongo-driver v1.10.0

## Test and build

Test the app:

```bash
make test
```

## Usage

Install dependencies:

```bash
go get github.com/javiertlopez/awesome
```

### Example

```go
package main

import (
	"net/http"
	"os"
	"time"

	"github.com/javiertlopez/awesome"

	"github.com/sirupsen/logrus"
)

const (
	writeTimeout = 15 * time.Second
	readTimeout  = 15 * time.Second
	idleTimeout  = 60 * time.Second
)

var (
	application awesome.App
	commit      string
	version     string
)

func main() {
	// Environment variables
	addr := os.Getenv("ADDR")
	mongoString := os.Getenv("MONGO_STRING")
	muxTokenID := os.Getenv("MUX_TOKEN_ID")
	muxTokenSecret := os.Getenv("MUX_TOKEN_SECRET")
	muxKeyID := os.Getenv("MUX_KEY_ID")
	muxKeySecret := os.Getenv("MUX_KEY_SECRET")

	// Create a logrus logger and set up the output format as JSON
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	application = awesome.New(
		awesome.AppConfig{
			Commit:         commit,
			Version:        version,
			MongoDB:        mongoDB,
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

	// Start the server
	logger.Fatal(srv.ListenAndServe())
}
```

## License

Licensed under [MIT License](LICENSE). © 2022 Hiram Torres Lopez.