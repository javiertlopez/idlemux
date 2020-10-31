# awesome [![Build Status](https://travis-ci.com/javiertlopez/awesome.svg?token=pyy6Hs7N6KLZpHXbFbXd&branch=main)](https://travis-ci.com/javiertlopez/awesome) [![codecov](https://codecov.io/gh/javiertlopez/awesome/branch/main/graph/badge.svg?token=I8D2Z4TZX4)](undefined)

## Overview
Awesome is a REST API that stores a Video Library in Mongo Atlas.

## Dependencies

- Go 1.15.3
- github.com/google/uuid v1.1.2
- github.com/gorilla/mux v1.8.0
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

Export environment environment variable with **Mongo Atlas Connection String**:

```bash
export MONGO_STRING=connectionString
```

Run app directly from terminal:

```bash
go run .
```

### From Visual Studio

The project is setup to run directly from Visual Studio Run tab. Just fill the credentials, and rename the file from `dev.env.example` to `dev.env`.