# awesome [![Build Status](https://travis-ci.com/javiertlopez/awesome.svg?token=pyy6Hs7N6KLZpHXbFbXd&branch=main)](https://travis-ci.com/javiertlopez/awesome) [![codecov](https://codecov.io/gh/javiertlopez/awesome/branch/main/graph/badge.svg?token=I8D2Z4TZX4)](undefined)

## Overview
Awesome is a REST API that stores a Video Library in Mongo Atlas and Mux.com for video ingestion.

## Dependencies

- Go 1.15.3
- github.com/google/uuid v1.1.2
- github.com/gorilla/mux v1.8.0
- github.com/muxinc/mux-go v0.9.0
- github.com/sirupsen/logrus v1.7.0
- go.mongodb.org/mongo-driver v1.4.2

For testing
- github.com/stretchr/testify v1.6.1

## Development

Install dependencies:

```bash
go get
```

### From Terminal

Export environment environment variable with **Mongo Atlas Connection String** and **Mux credentials**:

```bash
export MONGO_STRING=connectionString
export MUX_TOKEN_ID=muxTokenID
export MUX_TOKEN_SECRET=muxTokenSecret
```

Run app directly from terminal:

```bash
go run .
```

### From Visual Studio

The project is setup to run directly from Visual Studio Run tab. Just fill the credentials, and rename the file from `dev.env.example` to `dev.env`.

## Docker

Create a docker image

```bash
docker build -t javiertlopez/awesome-api .
```

Run the docker container

```bash
docker run --env-file=dev.env -d -p 8080:8080 javiertlopez/awesome-api
```

**Note.** It is required to fill the `dev.env` file.
