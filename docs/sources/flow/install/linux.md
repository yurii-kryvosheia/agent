---
title: Linux systems
weight: 300
---

# Install Grafana Agent Flow on Linux systems

Grafana Agent Flow can be installed as a systemd service on various AMD64 and
ARM64 Linux systems:

* [Debian-based systems](#install-on-debian-based-systems)
* [RedHat-based systems](#install-on-redhat-based-systems)

## Install on Debian-based systems

To install Grafana Agent Flow on Debian-based systems (such as Debian or
Ubuntu), complete the following steps:

1. Open a terminal and run the following command to install Grafana's package repository:

   ```shell
   mkdir -p /etc/apt/keyrings/
   wget -q -O - https://apt.grafana.com/gpg.key | gpg --dearmor > /etc/apt/keyrings/grafana.gpg
   echo "deb [signed-by=/etc/apt/keyrings/grafana.gpg] https://apt.grafana.com stable main" | tee /etc/apt/sources.list.d/grafana.list
   ```

2. With the repository installed, update the list of available packages:

   ```shell
   sudo apt-get update
   ```

3. Install the Grafana Agent Flow package:

   ```shell
   sudo apt-get install grafana-agent-flow
   ```

## Install on RedHat-based systems

To install Grafana Agent Flow on RedHat-based systems (such as CentOS, Fedora,
or RedHat Enterprise Linux), complete the following steps:

1. Create `/etc/yum.repos.d/grafana.repo` with the following content:

   ```
   [grafana]
   name=grafana
   baseurl=https://rpm.grafana.com
   repo_gpgcheck=1
   enabled=1
   gpgcheck=1
   gpgkey=https://rpm.grafana.com/gpg.key
   sslverify=1
   sslcacert=/etc/pki/tls/certs/ca-bundle.crt
   ```

2. Verify that the repository is properly configured using
   `yum-config-manager`:

   ```shell
   yum-config-manager grafana
   ```

3. Install the Grafana Agent Flow package:

   ```shell
   sudo yum install grafana-agent-flow
   ```

## Operation guide

After installing Grafana Agent Flow on Linux, it will be exposed as a
[systemd][] service.

[systemd]: https://systemd.io/

### Run Grafana Agent Flow

To run Grafana Agent Flow, run the following command in a terminal:

```shell
sudo systemctl start grafana-agent-flow
```

To check the status of Grafana Agent Flow, run the following command in a
terminal:

```shell
sudo systemctl status grafana-agent-flow
```

### Run Grafana Agent Flow on startup

To automatically run Grafana Agent Flow when the system starts, run the
following command in a terminal:

```shell
sudo systemctl enable grafana-agent-flow.service
```

### Configuring Grafana Agent Flow

To configure Grafana Agent Flow when installed on Linux, perform the following
steps:

1. Edit the default configuration file at `/etc/grafana-agent-flow.river`.

2. Run the following command in a terminal to reload the configuration file:

   ```shell
   sudo systemctl reload grafana-agent-flow
   ```

To change the configuration file used by the service, perform the following steps:

1. Edit the environment file for the service:

   * Debian-based systems: edit `/etc/default/grafana-agent-flow`
   * RedHat-based systems: edit `/etc/sysconfig/grafana-agent-flow`

2. Change the contents of the `CONFIG_FILE` environment variable to point to
   the new configuration file to use.

3. Restart the Grafana Agent Flow service:

   ```shell
   sudo systemctl restart grafana-agent-flow
   ```

### Passing additional command-line flags

By default, the Grafana Agent Flow service will launch with the [run][]
command, passing the following flags:

* `--storage.path=/var/lib/grafana-agent-flow`

To pass additional command-line flags to the Grafana Agent Flow binary, perform
the following steps:

1. Edit the environment file for the service:

   * Debian-based systems: edit `/etc/default/grafana-agent-flow`
   * RedHat-based systems: edit `/etc/sysconfig/grafana-agent-flow`

2. Change the contents of the `CUSTOM_ARGS` environment variable to specify
   command-line flags to pass.

3. Restart the Grafana Agent Flow service:

   ```shell
   sudo systemctl restart grafana-agent-flow
   ```

To see the list of valid command-line flags that can be passed to the service,
refer to the documentation for the [run][] command.

[run]: {{< relref "../reference/cli/run.md" >}}

### Viewing Grafana Agent Flow logs

Logs of Grafana Agent Flow can be found by running the following command in a
terminal:

```shell
sudo journalctl -u grafana-agent-flow
```

