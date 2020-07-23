<img src="assets/logo.png" alt="geoip-updater" width="128px" style="display: block; margin-left: auto; margin-right: auto"/>

<p align="center">
  <a href="https://github.com/crazy-max/geoip-updater/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/geoip-updater.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/geoip-updater/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/geoip-updater/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/geoip-updater/actions?workflow=build"><img src="https://img.shields.io/github/workflow/status/crazy-max/geoip-updater/build?label=build&logo=github&style=flat-square" alt="Build Status"></a>
  <a href="https://hub.docker.com/r/crazymax/geoip-updater/"><img src="https://img.shields.io/docker/stars/crazymax/geoip-updater.svg?style=flat-square&logo=docker" alt="Docker Stars"></a>
  <a href="https://hub.docker.com/r/crazymax/geoip-updater/"><img src="https://img.shields.io/docker/pulls/crazymax/geoip-updater.svg?style=flat-square&logo=docker" alt="Docker Pulls"></a>
  <br /><a href="https://goreportcard.com/report/github.com/crazy-max/geoip-updater"><img src="https://goreportcard.com/badge/github.com/crazy-max/geoip-updater?style=flat-square" alt="Go Report"></a>
  <a href="https://app.codacy.com/manual/crazy-max/geoip-updater"><img src="https://img.shields.io/codacy/grade/39bc1954d42b4b77b5efe7fe3c7b9a17.svg?style=flat-square" alt="Code Quality"></a>
  <a href="https://github.com/sponsors/crazy-max"><img src="https://img.shields.io/badge/sponsor-crazy--max-181717.svg?logo=github&style=flat-square" alt="Become a sponsor"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

---

## What is geoip-updater?

**geoip-updater** :globe_with_meridians: is a CLI application written in [Go](https://golang.org/) and delivered as a
[single executable]({{ config.repo_url }}releases/latest) (and a
[Docker image](https://hub.docker.com/r/crazymax/geoip-updater/)) that lets you download and update
[MaxMind](https://www.maxmind.com/)'s GeoIP2 databases on a time-based schedule.

With Go, this can be done with an independent binary distribution across all platforms and architectures that Go supports.
This support includes Linux, macOS, and Windows, on architectures like amd64, i386, ARM, PowerPC, and others.

## Features

* Support for MMDB and CSV databases
* List of Edition IDs currently supported are available [here](https://github.com/crazy-max/geoip-updater/blob/master/pkg/maxmind/editionid.go#L10-L18).
* Archive authenticity checked
* Internal cron implementation through go routines

## License

This project is licensed under the terms of the MIT license.
