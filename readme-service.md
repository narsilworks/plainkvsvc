# systemd Service Unit for Traefik

Adapted from [caddy systemd Service Unit](https://github.com/mholt/caddy/tree/e2635666730e24bfbc2408811be089502338cbc4/dist/init/linux-systemd)

The provided file should work with systemd version 219 or later. It might work with earlier versions.
The easiest way to check your systemd version is to run `systemctl --version`.

## Instructions

We will assume the following:

* that you want to run traefik as user `traefik` and group `traefik`, with UID and GID 321
* you are working from a non-root user account that can use 'sudo' to execute commands as root

Adjust as necessary or according to your preferences.

First, put the traefik binary in the system wide binary directory and give it
appropriate ownership and permissions:

```bash
sudo cp /path/to/traefik /usr/local/bin
sudo chown root:root /usr/local/bin/traefik
sudo chmod 755 /usr/local/bin/traefik
```

Give the traefik binary the ability to bind to privileged ports (e.g. 80, 443) as a non-root user:

```bash
sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/traefik
```

Set up the user, group, and directories that will be needed:

```bash
sudo groupadd -g 321 traefik
sudo useradd \
  -g traefik --no-user-group \
  --home-dir /var/www --no-create-home \
  --shell /usr/sbin/nologin \
  --system --uid 321 traefik

sudo mkdir /etc/traefik
sudo mkdir /etc/traefik/acme
sudo chown -R root:root /etc/traefik
sudo chown -R traefik:traefik /etc/traefik/acme
```

Place your traefik configuration file ("traefik.toml") in the proper directory
and give it appropriate ownership and permissions:

```bash
sudo cp /path/to/traefik.toml /etc/traefik/
sudo chown root:root /etc/traefik/traefik.toml
sudo chmod 644 /etc/traefik/traefik.toml
```

Install the systemd service unit configuration file, reload the systemd daemon,
and start traefik:

```bash
sudo cp /path/to/traefik.service /etc/systemd/system/
sudo chown root:root /etc/systemd/system/traefik.service
sudo chmod 644 /etc/systemd/system/traefik.service
sudo systemctl daemon-reload
sudo systemctl start traefik.service
```

Have the traefik service start automatically on boot if you like:

```bash
sudo systemctl enable traefik.service
```

If traefik doesn't seem to start properly you can view the log data to help figure out what the problem is:

```bash
journalctl --boot -u traefik.service
```

If your GNU/Linux distribution does not use *journald* with *systemd* then check any logfiles in `/var/log`.

If you want to follow the latest logs from traefik you can do so like this:

```bash
journalctl -f -u traefik.service
```