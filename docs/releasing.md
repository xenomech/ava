# Releasing

Three deployables, two release paths.

| What | Where | Trigger |
| --- | --- | --- |
| API server | Railway service `ava-api` | push to `main` |
| Web UI | Railway service `ava-web` | push to `main` |
| Hub binaries | GitHub Release | push a `v*` tag |

The hub is the odd one out: it runs on someone's Raspberry Pi, so it ships as a
downloadable binary rather than a deploy.

## Cutting a release

```sh
git tag v0.2.0
git push origin v0.2.0
```

That runs `.github/workflows/release.yml`:

1. **verify** — `make test` and `make check-boundary` across all three Go modules
2. **hub** — cross-compiles `avahub` and `avactl` for `linux/arm64`, `linux/amd64`
   and `darwin/arm64`, with the version stamped in
3. **publish** — SHA256 checksums, notes generated from commit subjects, GitHub Release

Nothing publishes if `verify` fails.

Versions are `vMAJOR.MINOR.PATCH`. Until the wire protocol between hub and API is
stable, treat any change to `pkg/wire` as a minor bump and expect hubs to need
updating alongside the server.

### Release notes

`.github/scripts/release-notes.sh` groups commit subjects by conventional-commit
prefix into Breaking / Features / Fixes / Other. Run it locally to preview:

```sh
.github/scripts/release-notes.sh v0.2.0
```

With no earlier tag it walks the entire history, so the first release will produce
a very long changelog. Pass an explicit starting point to trim it:

```sh
.github/scripts/release-notes.sh v0.2.0 v0.1.0
```

## Installing a hub on a Raspberry Pi

64-bit Raspberry Pi OS. The short way — installs the latest release, sets up the
`avahub` systemd service, and prints the pairing code:

```sh
curl -fsSL https://raw.githubusercontent.com/xenomech/ava/main/scripts/install-hub.sh \
  | sudo bash -s -- --api https://api.example.com --broker wss://broker.example.com:443
```

The broker URL must be encrypted. The hub authenticates with a password, so a
plaintext `tcp://` link to a public host hands that password to anyone on the
path — the hub refuses to connect over one. Two ways to encrypt it:

| Scheme | When |
|---|---|
| `wss://host:443` | A proxy in front of the broker terminates TLS. This is the Railway path, using its managed certificate — nothing to provision. |
| `ssl://host:8883` | The broker holds the certificate itself. Set `MQTT_TLS_CERTFILE` and `MQTT_TLS_KEYFILE` on the broker to open the listener. |

For a private certificate authority, point the hub at it with
`MQTT_CA_FILE=/etc/avahub/broker-ca.pem`. `tcp://` remains fine for a broker on
the same private network — the API reaches it that way over Railway's internal
network.

Re-running it with no flags upgrades the binary in place; configuration and the
pairing state survive, so the hub does not need re-pairing. `--name` overrides
the hub name and `--version vX.Y.Z` pins a release.

The long way, by hand:

```sh
curl -fsSLO https://github.com/xenomech/ava/releases/download/v0.2.0/avahub-linux-arm64
curl -fsSLO https://github.com/xenomech/ava/releases/download/v0.2.0/checksums.txt
sha256sum --check --ignore-missing checksums.txt

chmod +x avahub-linux-arm64
sudo install -m 0755 avahub-linux-arm64 /usr/local/bin/avahub
```

Check the version, then pair it:

```sh
avahub --version
avahub
```

The hub prints a pairing code; enter it in the web UI under Hubs. It writes
`avahub-state.json` next to its working directory — that file holds the hub's
tokens, so keep it at mode `0600` and never commit it.

To run it as a service, create `/etc/systemd/system/avahub.service`:

```ini
[Unit]
Description=Ava hub
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=/usr/local/bin/avahub
WorkingDirectory=/var/lib/avahub
Environment=API_BASE_URL=https://api.example.com/api/v1
Environment=MQTT_BROKER_URL=wss://broker.example.com:443
Restart=on-failure
RestartSec=5
User=avahub
UMask=0077

[Install]
WantedBy=multi-user.target
```

```sh
sudo systemctl enable --now avahub
```

Upgrading is `install` over the top plus `systemctl restart avahub`. The state file
survives, so the hub does not need re-pairing.

## Railway

Two services, both built from this repo with **Root Directory** set to `/` — the API
Dockerfile copies `backend/pkg`, so it needs the whole repo as build context.

Each service points at its own config file, set under
*Settings → Config-as-code → Railway Config File*:

| Service | Config file | Dockerfile |
| --- | --- | --- |
| `ava-api` | `railway.api.json` | `backend/services/api/Dockerfile` |
| `ava-web` | `railway.web.json` | `frontend/Dockerfile` |

`watchPatterns` in each config keeps a frontend-only commit from rebuilding the API
and the other way round.

You also need a Postgres service and an MQTT broker. Railway has a Postgres template;
there is no first-party MQTT one, so run Mosquitto as a service or point
`MQTT_BROKER_URL` at a hosted broker.

The broker refuses anonymous connections. Set `MQTT_USERNAME` and `MQTT_PASSWORD`
on both the broker and the API — the broker bootstraps that account on first start
and the API authenticates with it.

### Broker environment

The broker opens port 1883 in plaintext for the private network. Hubs connect
from outside it, so give them an encrypted listener as well:

| Variable | Notes |
| --- | --- |
| `MQTT_WEBSOCKET_PORT` | Set to `9001` and expose it as the service's HTTP port. Railway terminates TLS with its own certificate, so hubs use `wss://<broker-domain>:443` and there is nothing to provision. |
| `MQTT_TLS_CERTFILE` `MQTT_TLS_KEYFILE` | Open a native TLS listener instead, for a broker exposed directly. `MQTT_TLS_PORT` defaults to 8883 and `MQTT_TLS_CAFILE` adds a client-certificate authority. |

Never publish port 1883 itself. The credentials the hub sends over it are the
only thing standing between a stranger and someone's lights.

### `ava-api` environment

Required — the API will not work without these:

| Variable | Notes |
| --- | --- |
| `SERVER_ENV` | Use `production`. `local` selects `sslmode=disable`, which Railway Postgres refuses. |
| `JWT_SECRET` | No default, and the API refuses to start without it. At least 32 characters outside `local`: `openssl rand -hex 32`. |
| `DB_HOST` `DB_PORT` `DB_USER` `DB_PASSWORD` `DB_DATABASE` | Reference the Postgres service's variables. |
| `MQTT_BROKER_URL` | e.g. `tcp://mosquitto.railway.internal:1883`. Plaintext is allowed here only because the internal network never leaves Railway; a public host needs `wss://` or `ssl://`. |
| `APP_URL` | Public URL of the web service. Used in verification and invite emails. |
| `CORS_ALLOWED_ORIGINS` | The web service's public origin. The default is localhost, which blocks the deployed UI. |
| `COOKIE_DOMAIN` | Parent domain shared by web and API, so the auth cookies are sent. |
| `RESEND_API_KEY` | No default; email silently fails without it. |

`DB_SSL_MODE` overrides the mode `SERVER_ENV` would pick — leave it unset on Railway.
It exists so the production image can be run against a non-TLS Postgres locally:

```sh
docker run --network ava-course_ava-network \
  -e SERVER_ENV=production -e DB_SSL_MODE=disable \
  -e DB_HOST=postgres -e DB_PORT=5432 -e DB_USER=ava \
  -e DB_PASSWORD=ava_password -e DB_DATABASE=ava_db \
  -e MQTT_BROKER_URL=tcp://mosquitto:1883 -e JWT_SECRET=$(openssl rand -hex 32) \
  -p 8097:8000 ava-api:test
```

Optional, all defaulted: `PORT`, `LOG_LEVEL`, `CORS_ALLOWED_METHODS`,
`CORS_ALLOWED_HEADERS`, `CORS_MAX_AGE`, `JWT_ACCESS_EXPIRY`, `JWT_REFRESH_EXPIRY`,
`HUB_CODE_EXPIRY`, `HUB_POLL_INTERVAL`, `HUB_TOKEN_EXPIRY`, `RESEND_FROM_EMAIL`,
`RESEND_FROM_NAME`.

Railway injects `PORT`; the Dockerfile's healthcheck hardcodes 8000, so leave `PORT`
unset or set it to `8000`.

### `ava-web` environment

One variable, and it behaves differently from every other one here:

| Variable | Notes |
| --- | --- |
| `VITE_API_URL` | Full URL including `/api/v1`, e.g. `https://api.example.com/api/v1` |

**`VITE_API_URL` is baked in at build time, not read at runtime.** Vite substitutes it
into the bundle during `pnpm build`, so changing it in the Railway dashboard does
nothing until the service rebuilds. After changing it, trigger a redeploy.

Railway passes service variables to Docker builds as build args, which is why
`frontend/Dockerfile` declares `ARG VITE_API_URL` and fails the build when it is
empty — a silent fallback to `http://localhost:8000/api/v1` would produce a deployed
UI that quietly talks to nothing.

### How the web app is served

A Caddy container, not `vite preview` (a dev server) and not Railpack's static
detection. The SPA needs two things that are easy to get wrong:

- **History fallback** — TanStack Router owns the URL, so unknown paths must return
  `index.html` rather than 404.
- **Cache headers** — `/assets/*` is content-hashed and immutable, but `sw.js`,
  `index.html` and the manifest must never be long-cached. A cached service worker
  pins users to an old build; that has already happened once in this project.

`frontend/Caddyfile` encodes both.

## Config-as-code is on a deadline

Railway deprecated `railway.json` in favour of Infrastructure as Code
(`.railway/railway.ts`). **Existing config files stop being read on 2026-12-01.**

We are on the deprecated format deliberately: the IaC reference does not yet document
Dockerfile builders, watch patterns or restart policies, all of which these services
need. When it does, migrate with:

```sh
railway config migrate           # preview the generated .railway/railway.ts
railway config migrate --apply   # write it and clear the Railway Config File setting
```

Run it once per service.

## Broker authentication

The broker runs Mosquitto's dynamic-security plugin with `allow_anonymous false`.
Two kinds of account exist.

**The API** holds one long-lived account, named by `MQTT_USERNAME`. Its password is
`MQTT_PASSWORD`, and `docker/mosquitto/entrypoint.sh` creates the account the first
time the broker starts against an empty config volume. The account carries the dynsec
admin role plus an `ava-control-plane` role the API creates for itself at boot, which
grants publish on `ava/+/+/cmd` and `ava/+/+/apply`. Without that second role the API
can administer the broker but cannot send a single command — the admin role the
bootstrap creates deliberately excludes publishing to ordinary topics.

**Each hub** gets its own account, created by the API when the hub pairs and returned
in the token response as `broker.username` and `broker.password`. The hub stores them
in `avahub-state.json` at mode `0600` and uses them to connect. Re-pairing or a token
refresh rotates them. Revoking a hub deletes its account.

A hub's role permits exactly four things:

| | |
| --- | --- |
| publish | `ava/{tenant}/{hub}/state` |
| publish | `ava/{tenant}/{hub}/status` |
| subscribe | `ava/{tenant}/{hub}/cmd` |
| subscribe | `ava/{tenant}/{hub}/apply` |

Everything else is refused, including publishing to its own `cmd` topic, reading any
other hub's topics, subscribing to a wildcard, and writing to the dynsec control topic.
One customer's hub cannot see or drive another customer's devices even if its
credentials leak.

The credentials live only in the hub's state file and the broker's
`dynamic-security.json`, which sits on the `mosquitto_config` volume at mode `0600`.
The API never stores a hub's broker password — it generates one, hands it over once,
and forgets it. Losing the state file means re-pairing.

To rotate the API account, change `MQTT_PASSWORD` and run
`mosquitto_ctrl dynsec setClientPassword` against the broker; the bootstrap only runs
when `dynamic-security.json` does not already exist.
