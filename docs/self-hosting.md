# Self-hosting Ava

Running your own Ava means three things: the API with a Postgres database, an MQTT
broker, and the web app. Then a hub on a Raspberry Pi in the house with the bulbs.

Ava is early software. Treat it accordingly: back the database up, and expect to
update hubs when you update the server.

## What you need

- Somewhere to run a container for the API, and a Postgres 17 database.
- An MQTT broker. The repo ships a Mosquitto image configured for this
  (`docker/mosquitto/`); a hosted broker works too.
- Static hosting for the web app.
- A Raspberry Pi on 64-bit Raspberry Pi OS for the hub.
- Optionally a [Resend](https://resend.com) key — without it, verification and invite
  emails silently fail.

## Configuring the API

Copy `backend/services/api/.env.example` and set at least these:

| Variable | Notes |
| --- | --- |
| `SERVER_ENV` | `production` in a deployment. `local` selects `sslmode=disable`, which most managed Postgres refuses. |
| `JWT_SECRET` | No default. Generate one: `openssl rand -base64 48`. |
| `DB_HOST` `DB_PORT` `DB_USER` `DB_PASSWORD` `DB_DATABASE` | Your Postgres. |
| `MQTT_BROKER_URL` | e.g. `tcp://broker.internal:1883` |
| `MQTT_USERNAME` `MQTT_PASSWORD` | The API's own broker account. Must match the broker. |
| `APP_URL` | Public URL of the web app. Used in verification and invite emails. |
| `CORS_ALLOWED_ORIGINS` | The web app's origin. The default is localhost, which blocks a deployed UI. |
| `COOKIE_DOMAIN` | Leave empty for host-only cookies. Set it (`.example.com`) only when the API and web app sit on different subdomains. |
| `RESEND_API_KEY` | Email fails quietly without it. |

`DB_SSL_MODE` overrides whatever `SERVER_ENV` would pick, which is useful for running
a production image against a local non-TLS Postgres.

The schema is applied on boot; there is no separate migration step.

## Configuring the web app

One variable, and it behaves unlike the others:

```
VITE_API_URL=https://api.example.com/api/v1
```

**It is baked in at build time, not read at runtime.** Vite substitutes it into the
bundle during `pnpm build`, so changing it in a dashboard does nothing until the app
is rebuilt. `frontend/Dockerfile` fails the build when it is empty rather than
falling back to localhost, which would produce a deployed UI that talks to nothing.

Serve the built output as a single-page app: unknown paths must return `index.html`,
`/assets/*` is content-hashed and can be cached forever, and `sw.js`, `index.html`
and the manifest must never be long-cached — a cached service worker pins people to
an old build. `frontend/Caddyfile` encodes all of that if you want a starting point.

If you serve the API from the same origin as the web app, session cookies work
without any `COOKIE_DOMAIN` juggling. That is what the shipped Caddyfile does.

## Broker authentication

The broker runs Mosquitto's dynamic-security plugin with `allow_anonymous false`.

**The API** holds one long-lived account named by `MQTT_USERNAME`.
`docker/mosquitto/entrypoint.sh` creates it the first time the broker starts against
an empty config volume, using `MQTT_ADMIN_USERNAME` / `MQTT_ADMIN_PASSWORD`. The
account gets the dynsec admin role plus an `ava-control-plane` role the API creates
for itself at boot, which grants publish on `ava/+/+/cmd` and `ava/+/+/apply`.
Without that second role the API can administer the broker but cannot send a single
command.

**Each hub** gets its own account, created when the hub pairs and returned once in
the token response. The hub stores it in its state file at mode `0600`. A hub's role
permits exactly four things:

| | |
| --- | --- |
| publish | `ava/{tenant}/{hub}/state` |
| publish | `ava/{tenant}/{hub}/status` |
| subscribe | `ava/{tenant}/{hub}/cmd` |
| subscribe | `ava/{tenant}/{hub}/apply` |

Everything else is refused, including publishing to its own `cmd` topic and reading
another hub's topics. One home's hub cannot see or drive another home's devices even
if its credentials leak.

The API never stores a hub's broker password — it generates one, hands it over once,
and forgets it. Losing a hub's state file means re-pairing it.

## Installing a hub on a Raspberry Pi

64-bit Raspberry Pi OS. On the Pi:

```bash
curl -fsSL https://raw.githubusercontent.com/xenomech/ava/main/scripts/install-hub.sh \
  | sudo bash -s -- --api https://api.example.com --broker tcp://broker.example.com:1883
```

That verifies checksums, installs `avahub` and `avactl` to `/usr/local/bin`, creates
an `avahub` system user, writes `/etc/avahub/avahub.env`, starts the `avahub` systemd
service, and prints the pairing code. Enter the code in the web app under
Settings → Hubs.

Options: `--name "Living room hub"` (defaults to the hostname), `--version vX.Y.Z` to
pin a release instead of taking the latest.

Re-running with no flags upgrades in place; configuration and pairing state survive,
so the hub does not need re-pairing:

```bash
curl -fsSL https://raw.githubusercontent.com/xenomech/ava/main/scripts/install-hub.sh | sudo bash
```

Check on it with `systemctl status avahub` and `journalctl -u avahub -f`. The hub's
tokens live in `/var/lib/avahub/avahub-state.json` at mode `0600` — that file is the
hub's identity, so keep it out of backups you share.

### Doing it by hand

```bash
curl -fsSLO https://github.com/xenomech/ava/releases/download/v0.1.0/avahub-linux-arm64
curl -fsSLO https://github.com/xenomech/ava/releases/download/v0.1.0/checksums.txt
sha256sum --check --ignore-missing checksums.txt
chmod +x avahub-linux-arm64
sudo install -m 0755 avahub-linux-arm64 /usr/local/bin/avahub
```

Then run `avahub` with `API_BASE_URL` and `MQTT_BROKER_URL` set, and pair it.

## Things that catch people out

- **The hub speaks plain MQTT over TCP.** There is no websocket fallback, so if your
  broker is behind a platform that only exposes HTTP, the Pi needs a real TCP
  endpoint to connect to.
- **Two hubs on one network fight.** Both answer for the same bulbs and each reports
  the other's devices as missing. Run one hub per network for now.
- **`VITE_API_URL` needs a rebuild**, not a restart. See above.
- **32-bit Raspberry Pi OS will not work.** Releases are 64-bit only; the installer
  refuses rather than installing something that cannot run.
