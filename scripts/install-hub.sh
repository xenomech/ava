#!/usr/bin/env bash
# Install or upgrade the Ava hub on a Raspberry Pi; re-running with no flags upgrades in place without re-pairing.
# Usage: curl -fsSL https://raw.githubusercontent.com/xenomech/ava/main/scripts/install-hub.sh | sudo bash -s -- --api <url> --broker <url> [--name <name>] [--version vX.Y.Z]

set -euo pipefail

REPO="xenomech/ava"
BIN_DIR="/usr/local/bin"
STATE_DIR="/var/lib/avahub"
ENV_FILE="/etc/avahub/avahub.env"
UNIT_FILE="/etc/systemd/system/avahub.service"
SERVICE_USER="avahub"

say()  { printf '\033[1;32m==>\033[0m %s\n' "$*"; }
fail() { printf '\033[1;31merror:\033[0m %s\n' "$*" >&2; exit 1; }

api_url="${AVA_API_URL:-}"
broker_url="${AVA_BROKER_URL:-}"
hub_name="${AVA_HUB_NAME:-}"
version=""

while [ $# -gt 0 ]; do
  case "$1" in
    --api)     api_url="${2:-}"; shift 2 ;;
    --broker)  broker_url="${2:-}"; shift 2 ;;
    --name)    hub_name="${2:-}"; shift 2 ;;
    --version) version="${2:-}"; shift 2 ;;
    *) fail "unknown flag: $1" ;;
  esac
done

[ "$(id -u)" -eq 0 ] || fail "run me with sudo: curl ... | sudo bash -s -- --api ..."
command -v systemctl >/dev/null || fail "systemd is required (is this Raspberry Pi OS?)"
command -v curl >/dev/null || fail "curl is required"

# Releases are 64-bit only; 32-bit Raspberry Pi OS reports armv7l even on 64-bit hardware.
case "$(uname -m)" in
  aarch64 | arm64) arch="arm64" ;;
  x86_64)          arch="amd64" ;;
  armv7l | armv6l) fail "32-bit OS detected ($(uname -m)); the hub needs 64-bit Raspberry Pi OS" ;;
  *)               fail "unsupported architecture: $(uname -m)" ;;
esac

# ── Resolve the release ──────────────────────────────────────────────────────

if [ -z "$version" ]; then
  version=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null |
    sed -n 's/.*"tag_name": *"\(v[^"]*\)".*/\1/p' | head -1) || true
  [ -n "$version" ] || fail "no published release found for ${REPO} — cut one with 'git tag vX.Y.Z && git push origin vX.Y.Z' (see docs/releasing.md)"
fi

asset="avahub-linux-${arch}"
base="https://github.com/${REPO}/releases/download/${version}"

say "Installing avahub ${version} (${arch})"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

curl -fsSL -o "${tmp}/${asset}" "${base}/${asset}" ||
  fail "release ${version} has no ${asset} asset"
curl -fsSL -o "${tmp}/checksums.txt" "${base}/checksums.txt" ||
  fail "release ${version} has no checksums.txt"

(cd "$tmp" && sha256sum --check --ignore-missing --quiet checksums.txt) ||
  fail "checksum mismatch for ${asset} — refusing to install"

install -m 0755 "${tmp}/${asset}" "${BIN_DIR}/avahub"

# avactl ships in the same release; useful for on-Pi debugging, optional.
if curl -fsSL -o "${tmp}/avactl-linux-${arch}" "${base}/avactl-linux-${arch}" 2>/dev/null; then
  if (cd "$tmp" && sha256sum --check --ignore-missing --quiet checksums.txt); then
    install -m 0755 "${tmp}/avactl-linux-${arch}" "${BIN_DIR}/avactl"
  fi
fi

# ── User, state, configuration ───────────────────────────────────────────────

id "$SERVICE_USER" >/dev/null 2>&1 ||
  useradd --system --home-dir "$STATE_DIR" --shell /usr/sbin/nologin "$SERVICE_USER"

# The state file holds the hub's tokens; only the service user may read it, and upgrades never touch it.
install -d -m 0700 -o "$SERVICE_USER" -g "$SERVICE_USER" "$STATE_DIR"
install -d -m 0755 "$(dirname "$ENV_FILE")"

if [ -n "$api_url" ] || [ -n "$broker_url" ] || [ ! -f "$ENV_FILE" ]; then
  [ -n "$api_url" ] || fail "--api is required on first install (or when reconfiguring)"
  [ -n "$broker_url" ] || fail "--broker is required on first install (or when reconfiguring)"

  case "$api_url" in
    */api/v1) ;;
    */) api_url="${api_url}api/v1" ;;
    *)  api_url="${api_url}/api/v1" ;;
  esac

  {
    echo "API_BASE_URL=${api_url}"
    echo "MQTT_BROKER_URL=${broker_url}"
    echo "STATE_FILE=${STATE_DIR}/avahub-state.json"
    if [ -n "$hub_name" ]; then echo "HUB_NAME=${hub_name}"; fi
  } >"$ENV_FILE"
  chmod 0600 "$ENV_FILE"
  chown "$SERVICE_USER":"$SERVICE_USER" "$ENV_FILE"
  say "Wrote ${ENV_FILE}"
else
  say "Keeping existing ${ENV_FILE}"
fi

# ── Service ──────────────────────────────────────────────────────────────────

cat >"$UNIT_FILE" <<UNIT
[Unit]
Description=Ava hub
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=${BIN_DIR}/avahub
WorkingDirectory=${STATE_DIR}
EnvironmentFile=${ENV_FILE}
Restart=on-failure
RestartSec=5
User=${SERVICE_USER}
UMask=0077

[Install]
WantedBy=multi-user.target
UNIT

systemctl daemon-reload
systemctl enable --now avahub >/dev/null 2>&1
systemctl restart avahub

say "avahub ${version} is running"

# ── Pairing: surface the code, or already-paired status, from the journal ────

say "Waiting for the hub to report in..."

for _ in $(seq 1 30); do
  log=$(journalctl -u avahub --since "-2 min" --no-pager -o cat 2>/dev/null || true)

  code=$(printf '%s\n' "$log" | sed -n 's/.*"user_code": *"\([^"]*\)".*/\1/p' | tail -1)
  if [ -n "$code" ]; then
    uri=$(printf '%s\n' "$log" | sed -n 's/.*"verification_uri": *"\([^"]*\)".*/\1/p' | tail -1)
    printf '\n  \033[1mPairing code: %s\033[0m\n' "$code"
    [ -n "$uri" ] && printf '  Enter it at: %s\n' "$uri"
    printf '\n'
    exit 0
  fi

  if printf '%s' "$log" | grep -q "HUB_PAIRED"; then
    say "Hub is already paired — no code needed"
    exit 0
  fi

  sleep 1
done

say "No pairing code yet — follow the logs with: journalctl -u avahub -f"
