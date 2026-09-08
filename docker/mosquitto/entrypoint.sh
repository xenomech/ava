#!/bin/sh
set -eu

# Under /mosquitto/data, the one writable directory a volume can mount over without hiding mosquitto.conf.
config=/mosquitto/data/dynamic-security.json
conf_d=/mosquitto/data/conf.d

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

# Rebuilt every boot so a listener disappears when its configuration does, rather than lingering in the volume.
rm -rf "$conf_d"
mkdir -p "$conf_d"

# Websockets for a platform that terminates TLS at its edge, which is where Railway's managed certificate applies.
if [ -n "${MQTT_WEBSOCKET_PORT:-}" ]; then
  {
    echo "listener ${MQTT_WEBSOCKET_PORT}"
    echo "protocol websockets"
  } >"${conf_d}/websockets.conf"
  echo "websockets listener on ${MQTT_WEBSOCKET_PORT}"
fi

# Native MQTT over TLS, for a broker exposed directly instead of sitting behind a proxy.
if [ -n "${MQTT_TLS_CERTFILE:-}" ] && [ -n "${MQTT_TLS_KEYFILE:-}" ]; then
  for file in "$MQTT_TLS_CERTFILE" "$MQTT_TLS_KEYFILE"; do
    if [ ! -r "$file" ]; then
      echo "TLS is configured but $file is not readable" >&2
      exit 1
    fi
  done

  {
    echo "listener ${MQTT_TLS_PORT:-8883}"
    echo "certfile ${MQTT_TLS_CERTFILE}"
    echo "keyfile ${MQTT_TLS_KEYFILE}"
    if [ -n "${MQTT_TLS_CAFILE:-}" ]; then
      echo "cafile ${MQTT_TLS_CAFILE}"
    fi
  } >"${conf_d}/tls.conf"
  echo "TLS listener on ${MQTT_TLS_PORT:-8883}"
fi

# We run as root but mosquitto drops privileges, and a root-owned 0600 config silently denies every client.
chown -R mosquitto:mosquitto /mosquitto/data

exec /usr/sbin/mosquitto -c /mosquitto/config/mosquitto.conf
