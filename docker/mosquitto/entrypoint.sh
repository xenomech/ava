#!/bin/sh
set -eu

# Under /mosquitto/data, not /mosquitto/config: this is the one directory the
# broker writes to, and it is the one a persistent volume can be mounted over
# without hiding mosquitto.conf.
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

# This script runs as root so it can write to a freshly attached volume, but
# mosquitto drops to its own user before reading anything. Left root-owned at
# 0600 the broker cannot read its own security config, and it says so only as
# "File is not readable" while carrying on with an empty default — which denies
# every client instead of failing outright.
chown -R mosquitto:mosquitto /mosquitto/data

exec /usr/sbin/mosquitto -c /mosquitto/config/mosquitto.conf
