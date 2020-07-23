# Examples

In this section we quickly go over basic ways to run geoip-updater.

## Single database

Download `GeoLite2-City` database with licence key `0123456789ABCD`:

```shell
$ geoip_updater \
    --edition-ids GeoLite2-City \
    --license-key 0123456789ABCD
```

## Multi databases

Download `GeoLite2-City` and `GeoLite2-Country` databases with licence key `0123456789ABCD` to
`/usr/local/share/geoip` with log level to `debug` on a time-based schedule (`every week`)

```shell
$ geoip_updater \
    --edition-ids GeoLite2-City,GeoLite2-Country \
    --license-key 0123456789ABCD \
    --download-path /usr/local/share/geoip \
    --schedule "0 0 * * 0" \
    --log-level debug
```
