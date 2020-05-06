# Usage

`geoip-updater --edition-ids=GeoLite2-City,GeoLite2-Country,... --license-key=0123456789`

* `--help`: Show help text and exit.
* `--version`: Show version and exit.
* `--edition-ids`: Edition IDs list (comma separated) of MaxMind's GeoIP2 databases to fdownload. Currently supported edition IDs by geoip-updater are available [here](https://github.com/crazy-max/geoip-updater/blob/master/pkg/maxmind/editionid.go#L10-L18). **Required**.
* `--license-key`: [MaxMind License Key](prerequisites.md#license-key) in order to download databases. **Required**.
* `--download-path`: Directory where databases will be stored. (default to geoip-updater root directory).
* `--schedule <cron expression>`: [CRON expression](https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format) to schedule geoip-updater. (eg. `0 0 * * 0`).
* `--timezone <timezone>`: Timezone assigned to geoip-updater. (default `UTC`).
* `--log-level <level>`: Log level output. (default `info`).
* `--log-json`: Enable JSON logging output. (default `false`).
* `--log-caller`: Add file:line of the caller to log output. (default `false`).

## Example

```
geoip_updater \
  --license-key 0123456789ABCD \
  --edition-ids GeoLite2-City
```
> Download `GeoLite2-City` database with licence key `0123456789ABCD` to geoip-updater root directory

```
geoip_updater \
  --license-key 0123456789ABCD \
  --download-path /usr/local/share/geoip \
  --schedule "0 0 * * 0" \
  --log-level debug \
  --edition-ids GeoLite2-City,GeoLite2-Country
```
> Download `GeoLite2-City` and `GeoLite2-Country` databases with licence key `0123456789ABCD` to `/usr/local/share/geoip` with log level to `debug` on a time-based schedule (`every week`)
