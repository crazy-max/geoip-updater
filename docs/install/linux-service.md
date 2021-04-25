# Run as service on Debian based distro

## Using systemd

!!! warning
    Make sure to follow the instructions to [install from binary](binary.md) before.

To create a new service, paste this content in `/etc/systemd/system/geoip-updater.service`:

```
[Unit]
Description=geoip-updater
Documentation={{ config.site_url }}
After=syslog.target
After=network.target

[Service]
RestartSec=2s
Type=simple
User=geoip-updater
Group=geoip-updater
ExecStart=/usr/local/bin/geoip-updater --log-level info
Restart=always
#Environment=TZ=Europe/Paris
Environment=EDITION_IDS=GeoLite2-ASN,GeoLite2-City,GeoLite2-Country
Environment=LICENSE_KEY=0123456789ABCD
Environment=DOWNLOAD_PATH=/usr/local/share/geoip
Environment=SCHEDULE=0 0 * * 0

[Install]
WantedBy=multi-user.target
```

Change the user, group, and other required startup values following your needs.

Enable and start geoip-updater at boot:

```shell
sudo systemctl enable geoip-updater
sudo systemctl start geoip-updater
```

To view logs:

```shell
journalctl -fu geoip-updater.service
```
