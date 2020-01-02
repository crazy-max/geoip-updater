# Installation from binary

## Download

geoip-updater binaries are available in [releases](https://github.com/crazy-max/geoip-updater/releases) page.

Choose the archive matching the destination platform and extract geoip-updater:

```
wget -qO- https://github.com/crazy-max/geoip-updater/releases/download/v0.2.1/geoip-updater_0.2.1_linux_x86_64.tar.gz | tar -zxvf - geoip-updater
```

## Test

After getting the binary, it can be tested with `./geoip-updater --help` or moved to a permanent location.

```
$ ./geoip-updater --help
usage: geoip-updater --license-key=0123456789 [<flags>] <edition-ids>

Download MaxMind's GeoIP2 databases on a time-based schedule. More info:
https://github.com/crazy-max/geoip-updater

Flags:
  --help                    Show context-sensitive help (also try --help-long
                            and --help-man).
  --license-key=0123456789  MaxMind License Key.
  --download-path=./        Directory where databases will be stored.
  --schedule=0 0 * * 0      CRON expression format.
  --timezone="UTC"          Timezone assigned to geoip-updater.
  --log-level="info"        Set log level.
  --log-json                Enable JSON logging output.
  --version                 Show application version.

Args:
  <edition-ids>  MaxMind Edition ID dbs to download (comma separated).
```

## Server configuration

Steps below are the recommended server configuration.

### Prepare environment

Create user to run geoip-updater (ex. `geoipupd`)

```
groupadd geoipupd
useradd -s /bin/false -d /bin/null -g geoipupd geoipupd
```

### Create required directory structure

```
mkdir -p /usr/local/share/geoip
chown geoipupd:geoipupd /usr/local/share/geoip
```

### Copy binary to global location

```
cp geoip-updater /usr/local/bin/geoip-updater
```

## Running geoip-updater

After the above steps, two options to run geoip-updater:

### 1. Creating a service file (recommended)

See how to create [Linux service](linux-service.md) to start geoip-updater automatically.

### 2. Running from command-line/terminal

```
/usr/local/bin/geoip_updater --license-key 0123456789ABCD --download-path /usr/local/share/geoip --schedule "0 0 * * 0" GeoLite2-City,GeoLite2-Country
```

## Updating to a new version

You can update to a new version of geoip-updater by stopping it, replacing the binary at `/usr/local/bin/geoip-updater` and restarting the instance.

If you have carried out the installation steps as described above, the binary should have the generic name `geoip-updater`. Do not change this, i.e. to include the version number.
