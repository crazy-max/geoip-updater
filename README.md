<p align="center"><a href="https://github.com/crazy-max/geoip-updater" target="_blank"><img height="128" src="https://raw.githubusercontent.com/crazy-max/geoip-updater/master/.res/geoip-updater.png"></a></p>

<p align="center">
  <a href="https://github.com/crazy-max/geoip-updater/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/geoip-updater.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/geoip-updater/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/geoip-updater/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/geoip-updater/actions"><img src="https://github.com/crazy-max/geoip-updater/workflows/build/badge.svg" alt="Build Status"></a>
  <a href="https://hub.docker.com/r/crazymax/geoip-updater/"><img src="https://img.shields.io/docker/stars/crazymax/geoip-updater.svg?style=flat-square" alt="Docker Stars"></a>
  <a href="https://hub.docker.com/r/crazymax/geoip-updater/"><img src="https://img.shields.io/docker/pulls/crazymax/geoip-updater.svg?style=flat-square" alt="Docker Pulls"></a>
  <br /><a href="https://goreportcard.com/report/github.com/crazy-max/geoip-updater"><img src="https://goreportcard.com/badge/github.com/crazy-max/geoip-updater?style=flat-square" alt="Go Report"></a>
  <a href="https://www.codacy.com/app/crazy-max/geoip-updater"><img src="https://img.shields.io/codacy/grade/39bc1954d42b4b77b5efe7fe3c7b9a17.svg?style=flat-square" alt="Code Quality"></a>
  <a href="https://github.com/sponsors/crazy-max"><img src="https://img.shields.io/badge/sponsor-crazy--max-181717.svg?logo=github&style=flat-square" alt="Become a sponsor"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

## About

**geoip-updater** :globe_with_meridians: is a CLI application written in [Go](https://golang.org/) that lets you download and update [MaxMind](https://www.maxmind.com/)'s GeoIP2 databases on a time-based schedule.

## Features

* Internal cron implementation through go routines
* Support for MMDB and CSV databases
* List of Edition IDs currently supported are available [here](https://github.com/crazy-max/geoip-updater/blob/master/pkg/maxmind/editionid.go#L10-L18).
* Archive authenticity checked
* Official [Docker image available](doc/install/docker.md)

## Documentation

* [Prerequisites](doc/prerequisites.md)
* Install
  * [With Docker](doc/install/docker.md)
  * [From binary](doc/install/binary.md)
  * [Linux service](doc/install/linux-service.md)
* [Usage](doc/usage.md)

## How can I help ?

All kinds of contributions are welcome :raised_hands:! The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon: You can also support this project by [**becoming a sponsor on GitHub**](https://github.com/sponsors/crazy-max) :clap: or by making a [Paypal donation](https://www.paypal.me/crazyws) to ensure this journey continues indefinitely! :rocket:

Thanks again for your support, it is much appreciated! :pray:

## License

MIT. See `LICENSE` for more details.
