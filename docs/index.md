<img src="assets/logo.png" alt="geoip-updater" width="128px" style="display: block; margin-left: auto; margin-right: auto"/>

<p align="center">
  <a href="https://github.com/crazy-max/geoip-updater/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/geoip-updater.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/geoip-updater/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/geoip-updater/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/geoip-updater/actions?workflow=build"><img src="https://img.shields.io/github/actions/workflow/status/crazy-max/geoip-updater/build.yml?branch=master&label=build&logo=github&style=flat-square" alt="Build Status"></a>
  <a href="https://hub.docker.com/r/crazymax/geoip-updater/"><img src="https://img.shields.io/docker/stars/crazymax/geoip-updater.svg?style=flat-square&logo=docker" alt="Docker Stars"></a>
  <a href="https://hub.docker.com/r/crazymax/geoip-updater/"><img src="https://img.shields.io/docker/pulls/crazymax/geoip-updater.svg?style=flat-square&logo=docker" alt="Docker Pulls"></a>
  <br /><a href="https://goreportcard.com/report/github.com/crazy-max/geoip-updater"><img src="https://goreportcard.com/badge/github.com/crazy-max/geoip-updater?style=flat-square" alt="Go Report"></a>
  <a href="https://github.com/sponsors/crazy-max"><img src="https://img.shields.io/badge/sponsor-crazy--max-181717.svg?logo=github&style=flat-square" alt="Become a sponsor"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

---

## What is geoip-updater?

**geoip-updater** :globe_with_meridians: keeps [MaxMind](https://www.maxmind.com/)'s GeoIP2
databases up to date. Configure the editions you need, provide your MaxMind
credentials, and let it download fresh MMDB or CSV archives on a schedule for
local use.

It is available as a [single executable]({{ config.repo_url }}releases/latest)
and a [container image](https://hub.docker.com/r/crazymax/geoip-updater/), so you
can run it on a host directly or as a containerized job in an existing stack.

## Features

* Downloads and refreshes supported GeoIP2 and GeoLite2 databases automatically
* Supports both MMDB and CSV database formats
* Runs on a configurable cron schedule without external scheduling tools
* Verifies downloaded archives before updating local database files
* Supports multiple Edition IDs in a single run

The full list of supported Edition IDs is available [here](https://github.com/crazy-max/geoip-updater/blob/master/pkg/maxmind/editionid.go).

## License

This project is licensed under the terms of the MIT license.
