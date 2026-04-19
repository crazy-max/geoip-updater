# Installation with Docker

## About

geoip-updater publishes automatically updated Docker :whale: images to multiple registries:

| Registry                                                                                                 | Image                             |
|----------------------------------------------------------------------------------------------------------|-----------------------------------|
| [Docker Hub](https://hub.docker.com/r/crazymax/geoip-updater/)                                           | `crazymax/geoip-updater`          |
| [GitHub Container Registry](https://github.com/users/crazy-max/packages/container/package/geoip-updater) | `ghcr.io/crazy-max/geoip-updater` |

You can use the latest stable tag directly or pair it with another service that
keeps Docker images up to date.

!!! note
    Want to be notified of new releases? Check out :bell: [Diun (Docker Image Update Notifier)](https://github.com/crazy-max/diun) project!

The image is available for the following platforms:

```shell
$ docker run --rm mplatform/mquery crazymax/geoip-updater:latest
Image: crazymax/geoip-updater:latest
 * Manifest List: Yes
 * Supported platforms:
   - linux/amd64
   - linux/arm/v6
   - linux/arm/v7
   - linux/arm64
   - linux/386
   - linux/ppc64le
```

This reference setup uses Docker Compose, but installing Docker Compose
itself is outside the scope of this documentation. Follow the official
[install instructions](https://docs.docker.com/compose/install/) if needed.

## Usage

```yaml
services:
  geoip-updater:
    image: crazymax/geoip-updater:latest
    container_name: geoip-updater
    volumes:
      - "./data:/data"
    environment:
      - "EDITION_IDS=GeoLite2-ASN,GeoLite2-City,GeoLite2-Country"
      - "LICENSE_KEY=0123456789ABCD"
      - "DOWNLOAD_PATH=/data"
      - "SCHEDULE=0 0 * * 0"
      - "LOG_LEVEL=info"
      - "LOG_JSON=false"
    restart: always
```

Edit this example to match your environment, then run the following commands to
start geoip-updater:

```shell
docker compose up -d
docker compose logs -f
```

Or run it directly with:

```shell
docker run -d --name geoip-updater \
  -e "EDITION_IDS=GeoLite2-ASN,GeoLite2-City,GeoLite2-Country" \
  -e "LICENSE_KEY=0123456789ABCD" \
  -e "DOWNLOAD_PATH=/data" \
  -e "SCHEDULE=0 0 * * 0" \
  -e "LOG_LEVEL=info" \
  -e "LOG_JSON=false" \
  -v "$(pwd)/data:/data" \
  crazymax/geoip-updater:latest
```

To upgrade to the latest release:

```shell
docker compose pull
docker compose up -d
```
