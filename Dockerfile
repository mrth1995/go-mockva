FROM golang:1.23-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN VERSION=$(cat version.json | grep -o '"version"[[:space:]]*:[[:space:]]*"[^"]*"' | cut -d'"' -f4) && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X github.com/mrth1995/go-mockva/pkg/version.Version=${VERSION}" -o /src/bin/app ./cmd

FROM alpine:3.20
RUN addgroup -S app && adduser -S -G app app
WORKDIR /srv
COPY --from=builder /src/bin/app /srv/app
COPY --from=builder /src/swagger-ui /srv/swagger-ui
COPY --from=builder /src/pkg/migration /srv/pkg/migration
RUN chown -R app:app /srv
USER app
EXPOSE 8080
ENTRYPOINT ["/srv/app"]

