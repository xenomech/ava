#!/bin/sh
set -eu

# Under /mosquitto/data, the one writable directory a volume can mount over without hiding mosquitto.conf.
config=/mosquitto/data/dynamic-security.json

mkdir -p /mosquitto/data

if [ ! -f "$config" ]; then
  if [ -z "${MQTT_ADMIN_PASSWORD:-}" ]; then
    echo "MQTT_ADMIN_PASSWORD is required to bootstrap the broker" >&2
    exit 1
  fi

  mosquitto_ctrl dynsec init "$config" "${MQTT_ADMIN_USERNAME:-ava-api}" "$MQTT_ADMIN_PASSWORD"
  chmod 0600 "$config"
  echo "bootstrapped dynamic security for ${MQTT_ADMIN_USERNAME:-ava-api}"
fi

# We run as root but mosquitto drops privileges, and a root-owned 0600 config silently denies every client.
chown -R mosquitto:mosquitto /mosquitto/data

exec /usr/sbin/mosquitto -c /mosquitto/config/mosquitto.conf
