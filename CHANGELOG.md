# Changelog

## 1.2.0 (2021/02/19)

* Refactor CI and dev workflow with buildx bake (#54)
  * Add `image-local` target
  * Single job for artifacts and image
  * Add `armv5`, `ppc64le` and `s390x` artifacts
  * Upload artifacts
  * Validate
* Remove `linux/s390x` Docker platform support for now
* Bump github.com/alecthomas/kong from 0.2.12 to 0.2.15 (#64)
* Bump github.com/stretchr/testify from 1.6.1 to 1.7.0 (#58)
* MkDocs Materials 6.2.8
* Bump github.com/mholt/archiver/v3 from 3.3.2 to 3.5.0 (#47)

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
