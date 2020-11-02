# Changelog

## 1.1.0 (2020/11/02)

* Use embedded tzdata package
* Remove `--timezone` flag
* Docker image also available on [GitHub Container Registry](https://github.com/users/crazy-max/packages/container/package/geoip-updater)
* Switch to Docker actions
* Go 1.15
* Update deps

## 1.0.1 (2020/05/06)

* Remove unexpected log output

## 1.0.0 (2020/05/06)

* Switch to kong command-line parser
* Edition IDs are now passed through the `--edition-ids` flag
* Add `--log-caller` flag
* Flag `--log-json` not handled
* Latest Go 1.13
* Switch to Open Container Specification labels as label-schema.org ones are deprecated
* Update deps

## 0.2.1 (2020/01/02)

* Fix incorrect checksum from archive

## 0.2.0 (2020/01/02)

* Improve extraction process and logging
* Move `.md5` to work directory
* Use a work directory to download archives and checksums

## 0.1.1 (2020/01/02)

* Check write permissions

## 0.1.0 (2020/01/02)

* Initial version
