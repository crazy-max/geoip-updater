# Changelog

## 1.9.0 (2023/12/17)

* Go 1.21 (#153)
* Bump github.com/alecthomas/kong from 0.7.1 to 0.8.1 (#154)
* Bump github.com/rs/zerolog from 1.28.0 to 1.31.0 (#143 #155)
* Bump github.com/stretchr/testify from 1.8.1 to 1.8.4 (#138 #146)
* Bump golang.org/x/sys from 0.3.0 to 0.11.0 (#142 #151)

## 1.8.0 (2022/12/31)

* Go 1.19 (#128)
* Alpine Linux 3.17 (#133)
* Enhance workflow (#129)
* Bump github.com/rs/zerolog from 1.27.0 to 1.28.0 (#125)
* Bump github.com/docker/go-units from 0.4.0 to 0.5.0 (#126)
* Bump github.com/stretchr/testify from 1.8.0 to 1.8.1 (#130)
* Bump github.com/alecthomas/kong from 0.6.1 to 0.7.1 (#132)
* Bump golang.org/x/sys to 0.3.0 (#134)

## 1.7.0 (2022/07/17)

* Go 1.18 (#120)
* Alpine Linux 3.16 (#122)
* MkDocs Material 8.3.9 (#123)
* Bump github.com/rs/zerolog from 1.26.1 to 1.27.0 (#114)
* Bump github.com/stretchr/testify from 1.7.0 to 1.8.0 (#101 #119)
* Bump github.com/alecthomas/kong from 0.3.0 to 0.6.1 (#100 #116)

## 1.6.0 (2022/01/08)

* Alpine Linux 3.15
* Bump github.com/alecthomas/kong from 0.2.22 to 0.3.0 (#95)

## 1.5.0 (2021/12/26)

* Fix checksum URL (#94)
* Move from io/ioutil to io and os packages
* Move syscall to golang.org/x/sys
* Enhance dockerfiles (#93)
* Bump github.com/rs/zerolog from 1.25.0 to 1.26.0 (#86)
* Bump github.com/mholt/archiver/v3 from 3.5.0 to 3.5.1 (#85)
* Bump github.com/alecthomas/kong from 0.2.17 to 0.2.22 (#88 #91)
* Bump github.com/rs/zerolog from 1.24.0 to 1.26.1 (#84 #92)

## 1.4.0 (2021/09/05)

* Go 1.17 (#82)
* Add `linux/riscv64`, `darwin/arm64`, `windows/arm64` artifacts
* MkDocs Materials 7.2.6
* Wrong remaining time displayed
* Alpine Linux 3.14
* Bump github.com/ulikunitz/xz to v0.5.8
* Bump github.com/rs/zerolog from 1.21.0 to 1.24.0 (#78 #80 #83)
* Bump github.com/alecthomas/kong from 0.2.16 to 0.2.17 (#79)
* Bump codecov/codecov-action from 1 to 2

## 1.3.0 (2021/03/28)

* Bump github.com/rs/zerolog from 1.20.0 to 1.21.0 (#69)
* Docker meta v2 (#70)
* Deploy docs on workflow dispatch or tag
* Bump github.com/alecthomas/kong from 0.2.15 to 0.2.16 (#68)
* Go 1.16 (#67)
* Switch to goreleaser-xx (#66)

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
