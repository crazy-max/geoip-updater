# Installation from binary

## Download

geoip-updater binaries are available on [releases]({{ config.repo_url }}releases/latest) page.

Choose the archive matching the destination platform:

* [`geoip-updater_{{ git.tag | trim('v') }}_darwin_amd64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_darwin_amd64.tar.gz)
* [`geoip-updater_{{ git.tag | trim('v') }}_darwin_arm64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_darwin_arm64.tar.gz)
* [`geoip-updater_{{ git.tag | trim('v') }}_freebsd_386.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_freebsd_386.tar.gz)
* [`geoip-updater_{{ git.tag | trim('v') }}_freebsd_amd64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_freebsd_amd64.tar.gz)
* [`geoip-updater_{{ git.tag | trim('v') }}_linux_386.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_linux_386.tar.gz)
* [`geoip-updater_{{ git.tag | trim('v') }}_linux_amd64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_linux_amd64.tar
* [`geoip-updater_{{ git.tag | trim('v') }}_linux_arm64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_linux_arm64.tar.gz)
* [`geoip-updater_{{ git.tag | trim('v') }}_linux_armv5.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_linux_armv5.tar.gz)
* [`geoip-updater_{{ git.tag | trim('v') }}_linux_armv6.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_linux_armv6.tar.gz)
* [`geoip-updater_{{ git.tag | trim('v') }}_linux_armv7.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_linux_armv7.tar.gz)
* [`geoip-updater_{{ git.tag | trim('v') }}_linux_ppc64le.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_linux_ppc64le.tar.gz)
* [`geoip-updater_{{ git.tag | trim('v') }}_linux_riscv64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_linux_riscv64.tar.gz)
* [`geoip-updater_{{ git.tag | trim('v') }}_linux_s390x.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_linux_s390x.tar.gz).gz)
* [`geoip-updater_{{ git.tag | trim('v') }}_windows_386.zip`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_windows_386.zip)
* [`geoip-updater_{{ git.tag | trim('v') }}_windows_amd64.zip`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_windows_amd64.zip)

And extract geoip-updater:

```shell
wget -qO- {{ config.repo_url }}releases/download/v{{ git.tag | trim('v') }}/geoip-updater_{{ git.tag | trim('v') }}_linux_amd64.tar.gz | tar -zxvf - geoip-updater
```

After getting the binary, it can be tested with [`./geoip-updater --help`](../usage/cli.md) command and moved to a
permanent location.

## Server configuration

Steps below are the recommended server configuration.

### Prepare environment

Create user to run geoip-updater (ex. `geoip-updater`)

```shell
groupadd geoip-updater
useradd -s /bin/false -d /bin/null -g geoip-updater geoip-updater
```

### Create required directory structure

```shell
mkdir -p /usr/local/share/geoip
chown geoip-updater:geoip-updater /usr/local/share/geoip
```

### Copy binary to global location

```shell
cp geoip-updater /usr/local/bin/geoip-updater
```

## Running geoip-updater

After the above steps, two options to run geoip-updater:

### 1. Creating a service file (recommended)

See how to create [Linux service](linux-service.md) to start geoip-updater automatically.

### 2. Running from terminal

```shell
/usr/local/bin/geoip_updater \
  --edition-ids GeoLite2-City,GeoLite2-Country \
  --license-key 0123456789ABCD \
  --download-path /usr/local/share/geoip \
  --schedule "0 0 * * 0"
```

## Updating to a new version

You can update to a new version of geoip-updater by stopping it, replacing the binary at
`/usr/local/bin/geoip-updater` and restarting the instance.

If you have carried out the installation steps as described above, the binary should have the generic name
`geoip-updater`. Do not change this, i.e. to include the version number.
