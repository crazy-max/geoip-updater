# Usage

`geoip-updater --license-key=0123456789 [<flags>] <edition-ids>`

* `--help`: Show help text and exit. _Optional_.
* `--version`: Show version and exit. _Optional_.
* `--license-key`: [MaxMind License Key](prerequisites.md#license-key) in order to download databases. **Required**.
* `--download-path`: Folder where databases will be stored. (default to geoip-updater root folder).
* `--schedule <cron expression>` : [CRON expression](https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format) to schedule geoip-updater. _Optional_. (example: `0 0 * * 0`).
* `--timezone <timezone>` : Timezone assigned to geoip-updater. _Optional_. (default: `UTC`).
* `--log-level <level>` : Log level output. _Optional_. (default: `info`).
* `--log-json` : Enable JSON logging output. _Optional_. (default: `false`).

`<edition-ids>` is edition IDs list (comma separated) of the MaxMind's GeoIP2 databases you would like to update. Currently supported edition IDs by geoip-updater are available [here](https://github.com/crazy-max/geoip-updater/blob/master/pkg/maxmind/editionid.go#L10-L18). 

## Example

`geoip_updater --license-key 0123456789ABCD GeoLite2-City`
> Download `GeoLite2-City` database with licence key `0123456789ABCD` to geoip-updater root folder

`geoip_updater --license-key 0123456789ABCD --download-path /usr/local/share/geoip --schedule "0 0 * * 0" --log-level debug GeoLite2-City,GeoLite2-Country`
> Download `GeoLite2-City` and `GeoLite2-Country` databases with licence key `0123456789ABCD` to `/usr/local/share/geoip` with log level to `debug` on a time-based schedule (`every week`)
