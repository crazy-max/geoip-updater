# Examples

This section covers a few basic ways to run geoip-updater.

## Single database

Download the `GeoLite2-City` database with license key `0123456789ABCD`:

```shell
geoip-updater \
  --edition-ids GeoLite2-City \
  --license-key 0123456789ABCD
```

## Multiple databases

Download the `GeoLite2-City` and `GeoLite2-Country` databases with license key
`0123456789ABCD` to `/usr/local/share/geoip`, set the log level to `debug`, and
run the update on a weekly schedule:

```shell
geoip-updater \
  --edition-ids GeoLite2-City,GeoLite2-Country \
  --license-key 0123456789ABCD \
  --download-path /usr/local/share/geoip \
  --schedule "0 0 * * 0" \
  --log-level debug
```
