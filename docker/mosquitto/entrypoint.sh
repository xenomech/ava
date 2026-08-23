#!/bin/sh
set -eu

config=/mosquitto/config/dynamic-security.json

if [ ! -f "$config" ]; then
  if [ -z "${MQTT_ADMIN_PASSWORD:-}" ]; then
    echo "MQTT_ADMIN_PASSWORD is required to bootstrap the broker" >&2
    exit 1
  fi

  mosquitto_ctrl dynsec init "$config" "${MQTT_ADMIN_USERNAME:-ava-api}" "$MQTT_ADMIN_PASSWORD"
  chmod 0600 "$config"
  echo "bootstrapped dynamic security for ${MQTT_ADMIN_USERNAME:-ava-api}"
fi

exec /usr/sbin/mosquitto -c /mosquitto/config/mosquitto.conf
