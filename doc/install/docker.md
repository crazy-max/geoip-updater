# Installation with Docker

geoip-updater provides automatically updated Docker :whale: images within [Docker Hub](https://hub.docker.com/r/crazymax/geoip-updater). It is possible to always use the latest stable tag or to use another service that handles updating Docker images.

Following platforms for this image are available:

```
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

Environment variables can be used within your container:

* `TZ`: Timezone assigned to geoip-updater
* `EDITION_IDS`: Edition IDs list (comma separated) of the MaxMind's GeoIP2 databases to update (default `GeoLite2-ASN,GeoLite2-City,GeoLite2-Country`)
* `LICENSE_KEY`: [MaxMind License Key](../prerequisites.md#license-key) in order to download databases
* `DOWNLOAD_PATH`: Folder where databases will be stored (default `/data`)
* `SCHEDULE`: [CRON expression](https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format) to schedule geoip-updater
* `LOG_LEVEL`: Log level output (default `info`)
* `LOG_JSON`: Enable JSON logging output (default `false`)

Docker compose is the recommended way to run this image. Copy the content of folder [.res/compose](../../.res/compose) in `/opt/geoip-updater/` on your host for example. Edit the compose file with your preferences and run the following commands:

```
docker-compose up -d
docker-compose logs -f
```

Or use the following command :

```
docker run -d --name geoip-updater \
  -e "TZ=Europe/Paris" \
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

```
docker-compose pull
docker-compose up -d
```
