# parrot

Small Golang program to republish data from MQTT to MQTT.

## installation

### container

This will start a long-running container with mqtt cron in it which runs all
available subcommands. The container is available for `amd64` and `aarch64`.

```
podman run -e MQTT_HOST="tcp://hostname:1883" ghcr.io/petspalace/parrot:latest
```

*The `podman` command can be switched out `docker` if you wish.*

## usage

todo
