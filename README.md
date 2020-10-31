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
go get github.com/javiertlopez/awesome/api 
```

### From Terminal

Export environment environment variable with **Mongo Atlas Connection String** and **Mux credentials**:

```bash
export MONGO_STRING="connectionString"
export MUX_TOKEN_ID=muxTokenID
export MUX_TOKEN_SECRET=muxTokenSecret
```

Run app directly from terminal:

```bash
go run ./api
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

## Usage

The Awesome App runs at 8080 port.

**Endpoint:** `http://127.0.0.1:8080`

### Hello World

Hello World path: `/app/hello`

Example:

```bash
curl --location --request GET 'http://127.0.0.1:8080/app/hello'
```

Response:

```json
{"message":"Hello World!","status":200}
```

### Videos

POST Path: `/videos`

**Example 1.** Video without source file:

```bash
curl --location --request POST 'http://127.0.0.1:8080/videos' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Don'\''t Look Back In Anger",
    "description": "Oasis song from (What'\''s the Story) Morning Glory? album."
}'
```

Response:
```json
{
    "id": "81095f2d-76d9-4b7b-a3d8-727aa90617c3",
    "title": "Don't Look Back In Anger",
    "description": "Oasis song from (What's the Story) Morning Glory? album.",
    "asset": {},
    "created_at": "2020-10-31 20:05:40.127259884 +0000 UTC m=+317.213044791",
    "updated_at": "2020-10-31 20:05:40.127259884 +0000 UTC m=+317.213044791"
}
```

**Example 2.** Video with a source file:

Request:
```bash
curl --location --request POST 'http://127.0.0.1:8080/videos' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Wonderwall",
    "description": "Oasis song from (What'\''s the Story) Morning Glory? album.",
    "source_url": "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4"
}'
```

Response:
```json
{
    "id": "50035c43-c1f9-4ff3-beb2-783f145a8dfe",
    "title": "Wonderwall",
    "description": "Oasis song from (What's the Story) Morning Glory? album.",
    "asset": {
        "id": "zQJ4DsW01QbbBl2Rd3801bUgpWmQD01j5eUOytzjv8QF02c"
    },
    "created_at": "2020-10-31 20:07:08.759255341 +0000 UTC m=+405.845420795",
    "updated_at": "2020-10-31 20:07:08.759255341 +0000 UTC m=+405.845420795"
}
```


GET Path: `/videos/{id}`

**Example 1.** Video without source file:

Request:
```bash
curl --location --request GET 'http://127.0.0.1:8080/videos/81095f2d-76d9-4b7b-a3d8-727aa90617c3'
```

Response:
```json
{
    "id": "81095f2d-76d9-4b7b-a3d8-727aa90617c3",
    "title": "Don't Look Back In Anger",
    "description": "Oasis song from (What's the Story) Morning Glory? album.",
    "asset": {},
    "created_at": "2020-10-31 20:05:40.127 +0000 UTC",
    "updated_at": "2020-10-31 20:05:40.127 +0000 UTC"
}
```

**Example 2.** Video with a source file:

Request:
```bash
curl --location --request GET 'http://127.0.0.1:8080/videos/50035c43-c1f9-4ff3-beb2-783f145a8dfe'
```

Response:
```json
{
    "id": "50035c43-c1f9-4ff3-beb2-783f145a8dfe",
    "title": "Wonderwall",
    "description": "Oasis song from (What's the Story) Morning Glory? album.",
    "asset": {
        "id": "zQJ4DsW01QbbBl2Rd3801bUgpWmQD01j5eUOytzjv8QF02c",
        "created_at": "1604174828",
        "status": "ready",
        "duration": 23.857167,
        "max_stored_resolution": "HD",
        "max_stored_frame_rate": 29.97,
        "aspect_ratio": "16:9",
        "poster": "https://image.mux.com/7hDIO2xL101KjoSrQy018NtTI02XhXQOGLbJcHVLC00YAig/thumbnail.png?width=1920\u0026height=1080\u0026smart_crop=true\u0026time=7",
        "thumbnail": "https://image.mux.com/7hDIO2xL101KjoSrQy018NtTI02XhXQOGLbJcHVLC00YAig/thumbnail.png?width=640\u0026height=360\u0026smart_crop=true\u0026time=7",
        "sources": [
            {
                "id": "7hDIO2xL101KjoSrQy018NtTI02XhXQOGLbJcHVLC00YAig",
                "policy": "",
                "src": "https://stream.mux.com/7hDIO2xL101KjoSrQy018NtTI02XhXQOGLbJcHVLC00YAig.m3u8",
                "type": "application/x-mpegURL"
            }
        ]
    },
    "created_at": "2020-10-31 20:07:08.759 +0000 UTC",
    "updated_at": "2020-10-31 20:07:08.759 +0000 UTC"
}
```