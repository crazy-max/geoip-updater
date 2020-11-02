FROM --platform=${BUILDPLATFORM:-linux/amd64} tonistiigi/xx:golang AS xgo
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.15-alpine as builder

ARG VERSION=dev

ENV CGO_ENABLED 0
ENV GO111MODULE on
ENV GOPROXY https://goproxy.io,direct
COPY --from=xgo / /

RUN apk --update --no-cache add \
    build-base \
    gcc \
    git \
  && rm -rf /tmp/* /var/cache/apk/*

WORKDIR /app

COPY . ./
RUN go mod download

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
RUN go env
RUN go build -ldflags "-w -s -X 'main.version=${VERSION}'" -v -o geoip-updater cmd/main.go

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:latest

LABEL maintainer="CrazyMax"

ENV EDITION_IDS="GeoLite2-ASN,GeoLite2-City,GeoLite2-Country" \
  DOWNLOAD_PATH="/data"

RUN apk --update --no-cache add \
    ca-certificates \
    libressl \
  && rm -rf /tmp/* /var/cache/apk/*

COPY --from=builder /app/geoip-updater /usr/local/bin/geoip-updater
RUN geoip-updater --version

ENTRYPOINT [ "/usr/local/bin/geoip-updater" ]
