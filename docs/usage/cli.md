# Command Line

## Usage

```shell
geoip-updater [options]
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
      --log-level="info"          Set log level ($LOG_LEVEL).
      --log-json                  Enable JSON logging output ($LOG_JSON).
      --log-caller                Add file:line of the caller to log output
                                  ($LOG_CALLER).
```

## Environment variables

The following environment variables can be used in place:

| Name            | Default       | Description                                                                                                      |
|-----------------|---------------|------------------------------------------------------------------------------------------------------------------|
| `EDITION_IDS`   |               | Edition IDs list (comma separated) of MaxMind's GeoIP2 databases to download                                     |
| `LICENSE_KEY`   |               | [MaxMind License Key](prerequisites.md#license-key) in order to download databases                               |
| `DOWNLOAD_PATH` | _working dir_ | Directory where databases will be stored                                                                         |
| `SCHEDULE`      |               | [CRON expression](https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format) to schedule geoip-updater |
| `LOG_LEVEL`     | `info`        | Log level output                                                                                                 |
| `LOG_JSON`      | `false`       | Enable JSON logging output                                                                                       |
| `LOG_CALLER`    | `false`       | Enable to add `file:line` of the caller                                                                          |

!!! tip
    List of supported edition IDs available [here](../faq.md#supported-edition-ids).
