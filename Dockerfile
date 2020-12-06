# golang alpine 1.15.6
FROM golang:1.15.6-alpine as builder

ARG commit
ARG version

ENV commit=$commit
ENV version=$version

# SSL for HTTPS calls.
RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates

# Create appuser.
ENV USER=appuser
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 

RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR $GOPATH/src/github.com/javiertlopez/awesome/
COPY . .

RUN go get github.com/javiertlopez/awesome/cmd/container
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.commit=${commit} -X main.version=${version}" -o /go/bin/main ./cmd/container

# Small image
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable.
COPY --from=builder /go/bin/main /go/bin/main
# Use an unprivileged user.

USER appuser:appuser

# Port on which the service will be exposed.
EXPOSE 8080

ENTRYPOINT ["/go/bin/main"]