# Run as service on Debian based distro

## Using systemd

> :warning: Make sure to follow the instructions to [install from binary](binary.md) before.

Run the below command in a terminal:

```
sudo vim /etc/systemd/system/geoip-updater.service
```

Copy the sample [geoip-updater.service](../../.res/systemd/geoip-updater.service).

Change the user, group, and other required startup values following your needs.

Enable and start geoip-updater at boot:

```
sudo systemctl enable geoip-updater
sudo systemctl start geoip-updater
```

To view logs:

```
journalctl -fu geoip-updater.service
```
