# Installation with Docker

## About

geoip-updater provides automatically updated Docker :whale: images within several registries:

| Registry                                                                                         | Image                           |
|--------------------------------------------------------------------------------------------------|---------------------------------|
| [Docker Hub](https://hub.docker.com/r/crazymax/geoip-updater/)                             | `crazymax/geoip-updater`                 |
| [GitHub Container Registry](https://github.com/users/crazy-max/packages/container/package/geoip-updater)  | `ghcr.io/crazy-max/geoip-updater`        |

It is possible to always use the latest stable tag or to use another service that handles updating Docker images.

!!! note
    Want to be notified of new releases? Check out :bell: [Diun (Docker Image Update Notifier)](https://github.com/crazy-max/diun) project!

Following platforms for this image are available:

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
   - linux/s390x
```

This reference setup guides users through the setup based on `docker-compose`, but the installation of `docker-compose`
is out of scope of this documentation. To install `docker-compose` itself, follow the official
[install instructions](https://docs.docker.com/compose/install/).

## Usage

```yaml
version: "3.5"

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

Edit this example with your preferences and run the following commands to bring up geoip-updater:

```shell
$ docker-compose up -d
$ docker-compose logs -f
```

Or use the following command:

```shell
$ docker run -d --name geoip-updater \
    -e "EDITION_IDS=GeoLite2-ASN,GeoLite2-City,GeoLite2-Country" \
    -e "LICENSE_KEY=0123456789ABCD" \
    -e "DOWNLOAD_PATH=/data" \
    -e "SCHEDULE=0 0 * * 0" \
    -e "LOG_LEVEL=info" \
    -e "LOG_JSON=false" \
    -v "$(pwd)/data:/data" \
    crazymax/geoip-updater:latest
```

To upgrade your installation to the latest release:

```shell
$ docker-compose pull
$ docker-compose up -d
```
