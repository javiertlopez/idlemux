COMMIT			:=		$(shell git rev-parse --short HEAD)
VERSION			:=		$(shell git describe --always --long --dirty)
IMAGE_NAME		:=		javiertlopez/awesome
LAMBDA			:=		delivery

all: test

test:
	go test ./... -v

fmt:
	go fmt ./...

image:
	docker build --build-arg commit=${COMMIT} --build-arg version=${VERSION} -t ${IMAGE_NAME} .

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.commit=${COMMIT} -X main.version=${VERSION}" -o ./dist/main ./cmd/lambda
	chmod +x ./dist/main
	zip -r -j ./dist/function-${COMMIT}.zip ./dist/*

upload:
	aws lambda update-function-code \
		--function-name  $(LAMBDA) \
		--zip-file fileb://dist/function-${COMMIT}.zip

.PHONY: test fmt build upload