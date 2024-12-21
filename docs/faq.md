# FAQ

## Timezone

By default, all interpretations and scheduling are done with your local timezone
(`TZ` environment variable).

Cron schedule may also override the timezone to be interpreted in by providing
an additional space-separated field at the beginning of the cron spec, of the
form `CRON_TZ=<timezone>`:

```shell
geoip-updater --schedule "CRON_TZ=Asia/Tokyo */30 * * * *"
```

## Supported edition IDs

Here is the list of supported [edition IDs](usage/cli.md#options) by geoip-updater:

* `GeoIP2-City`
* `GeoIP2-City-CSV`
* `GeoIP2-Country`
* `GeoIP2-Country-CSV`
* `GeoLite2-ASN`
* `GeoLite2-ASN-CSV`
* `GeoLite2-City`
* `GeoLite2-City-CSV`
* `GeoLite2-Country`
* `GeoLite2-Country-CSV`
