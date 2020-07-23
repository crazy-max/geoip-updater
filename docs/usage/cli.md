# Command Line

## Usage

```shell
$ geoip-updater [options]
```

## Options

```
$ geoip-updater --help
Usage: geoip-updater --edition-ids=GeoLite2-City,GeoLite2-Country,... --license-key=0123456789

Download and update MaxMind's GeoIP2 databases on a time-based schedule. More
info: https://github.com/crazy-max/geoip-updater

Flags:
  -h, --help                      Show context-sensitive help.
      --version
      --edition-ids=GeoLite2-City,GeoLite2-Country,...
                                  MaxMind Edition ID dbs to download (comma
                                  separated) ($EDITION_IDS).
      --license-key=0123456789    MaxMind License Key ($LICENSE_KEY)
      --download-path=./          Directory where databases will be stored
                                  ($DOWNLOAD_PATH).
      --schedule=0 0 * * 0        CRON expression format ($SCHEDULE).
      --timezone="UTC"            Timezone assigned to geoip-updater ($TZ).
      --log-level="info"          Set log level ($LOG_LEVEL).
      --log-json                  Enable JSON logging output ($LOG_JSON).
      --log-caller                Add file:line of the caller to log output
                                  ($LOG_CALLER).
```

## Environment variables

Following environment variables can be used in place:

| Name               | Default       | Description   |
|--------------------|---------------|---------------|
| `EDITION_IDS`      |               | Edition IDs list (comma separated) of MaxMind's GeoIP2 databases to fdownload. Currently supported edition IDs by geoip-updater are available [here](https://github.com/crazy-max/geoip-updater/blob/master/pkg/maxmind/editionid.go#L10-L18) |
| `LICENSE_KEY`      |               | [MaxMind License Key](prerequisites.md#license-key) in order to download databases |
| `DOWNLOAD_PATH`    | _working dir_ | Directory where databases will be stored |
| `SCHEDULE`         |               | [CRON expression](https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format) to schedule geoip-updater |
| `TZ`               | `UTC`         | Timezone assigned |
| `LOG_LEVEL`        | `info`        | Log level output |
| `LOG_JSON`         | `false`       | Enable JSON logging output |
| `LOG_CALLER`       | `false`       | Enable to add `file:line` of the caller |
