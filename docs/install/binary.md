# Installation from binary

## Download

geoip-updater binaries are available on the [releases]({{ config.repo_url }}/releases/latest) page.

Choose the archive matching the destination platform:

* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_darwin_amd64.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_darwin_amd64.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_darwin_arm64.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_darwin_arm64.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_freebsd_386.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_freebsd_386.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_freebsd_amd64.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_freebsd_amd64.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_386.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_386.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_amd64.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_amd64.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_arm64.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_arm64.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_armv5.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_armv5.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_armv6.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_armv6.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_armv7.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_armv7.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_ppc64le.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_ppc64le.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_riscv64.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_riscv64.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_s390x.tar.gz`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_s390x.tar.gz)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_windows_386.zip`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_windows_386.zip)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_windows_amd64.zip`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_windows_amd64.zip)
* [`geoip-updater_{{ latest_stable_tag | trim('v') }}_windows_arm64.zip`]({{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_windows_arm64.zip)

Then extract `geoip-updater`:

```shell
wget -qO- {{ config.repo_url }}/releases/download/{{ latest_stable_tag }}/geoip-updater_{{ latest_stable_tag | trim('v') }}_linux_amd64.tar.gz | tar -zxvf - geoip-updater
```

After downloading the binary, test it with [`./geoip-updater --help`](../usage/cli.md)
and move it to a permanent location.

## Server configuration

The steps below describe the recommended server configuration.

### Prepare environment

Create a user to run geoip-updater, for example `geoip-updater`:

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

After the steps above, there are two ways to run geoip-updater:

### 1. Creating a service file (recommended)

See [Linux service](linux-service.md) to start geoip-updater automatically.

### 2. Running from terminal

```shell
/usr/local/bin/geoip-updater \
  --edition-ids GeoLite2-City,GeoLite2-Country \
  --license-key 0123456789ABCD \
  --download-path /usr/local/share/geoip \
  --schedule "0 0 * * 0"
```

## Updating to a new version

To update geoip-updater, stop it, replace the binary at
`/usr/local/bin/geoip-updater`, and restart the service or process.

If you followed the installation steps above, the binary should keep the generic
name `geoip-updater`. Do not rename it to include the version number.
