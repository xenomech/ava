# Contributing

Thanks for looking. This describes how to get Ava running locally, what the branches
mean, and what a reviewable change looks like.

Note there is no licence file yet, so contributions cannot be redistributed until one
is added. If you plan to send a substantial change, it is worth opening an issue first
to check that licensing has been settled.

## Setting up

Prerequisites: **Go 1.26+**, **pnpm 11+** on **Node 24** (what CI builds with),
**Docker**.

```bash
make setup      # air, golangci-lint, gofumpt, lefthook + git hooks
cp backend/services/api/.env.example backend/services/api/.env
cp frontend/apps/web/.env.example frontend/apps/web/.env
```

Set `MQTT_PASSWORD` in `backend/services/api/.env`; the broker will not start without
one. Then bring up the dependencies and the two dev servers:

```bash
docker compose --env-file backend/services/api/.env up -d postgres mosquitto
make dev                                  # API on :8000 with hot reload
pnpm -C frontend install
pnpm -C frontend dev:web                  # web on :3000
```

Compose publishes Postgres on `DB_PORT` and the broker on `MQTT_PORT`, and the API
connects to those same ports on localhost, so the two always agree.

`.env.example` sets `DB_PORT=5432`. If you already run Postgres locally, compose will
fail with `port is already allocated` — change `DB_PORT` to something free (55432,
say) and everything follows, because the API reads the same variable.

### Working on the hub

The hub is a separate Go module that runs on a Raspberry Pi, but it runs on your
machine too — it discovers real bulbs on whatever network you are on.

```bash
API_BASE_URL=http://localhost:8000/api/v1 MQTT_BROKER_URL=tcp://localhost:1893 make hub-run
```

It prints a pairing code; enter it in the web app under Settings → Hubs.
`avactl` is a companion CLI for poking at discovery without running the daemon.

## Branches

`stage` is **preproduction**. `main` is **production**. Work flows one way:

```
feature branch  ->  PR into stage  ->  verify in preproduction  ->  promote stage to main
```

Branch from `origin/stage`, and open pull requests against `stage`:

```bash
git fetch origin
git switch -c my-change origin/stage
gh pr create --base stage
```

Please do not open a feature PR against `main` — that ships to production without
passing through preproduction.

## Commits

Conventional commits, one short subject line, no body:

```
feat: aim the room switch with a scene carousel
fix: stop a stale sweep undoing a colour you just set
refactor: colocate each module and give the shell slots
```

`feat`, `fix`, `refactor`, `perf`, `style`, `test`, `docs`, `chore`. Release notes are
generated from these subjects, so write them for someone reading a changelog.

Keep each commit buildable on its own — the history gets bisected.

## Checks

Lefthook runs the fast checks on commit. Before opening a PR:

```bash
make lint            # golangci-lint across all three Go modules
make fmt             # gofumpt
make test            # Go tests (see the caveat below)
make check-boundary  # hub and api must not import each other

pnpm -C frontend check-types
pnpm -C frontend check          # oxlint + oxfmt
pnpm -C frontend build
```

### Running the Go tests properly

**A large part of the Go suite skips unless you give it a database.** The repository
and service tests call `t.Skip` when `AVA_TEST_DB_HOST` is unset, so `go test ./...`
reports `ok` while silently running almost nothing. CI does not set these either.

To actually run them, point at a scratch database — not your dev one, the tests write
freely:

```bash
docker compose --env-file backend/services/api/.env up -d postgres
docker exec ava-course-postgres-1 psql -U ava -d ava_db -c 'create database ava_test'

cd backend/services/api
AVA_TEST_DB_HOST=localhost \
AVA_TEST_DB_PORT=5432 \
AVA_TEST_DB_USER=ava \
AVA_TEST_DB_PASSWORD=ava_password \
AVA_TEST_DB_NAME=ava_test \
  go test ./... -count=1
```

Use whatever `DB_PORT` you set above. If you are unsure a test really ran, look for
`SKIP` in `go test -v` output — a skipped package still reports `ok`.

## Code conventions

**Comments are single line.** No block comments and no multi-line comments anywhere
in the repo — one line, saying the thing the code cannot.

**Frontend architecture is enforced by the linter.** `frontend/.oxlintrc.json` bans
deep imports between modules and keeps the layers pointing one way:

```
modules/  ->  shared/  ->  config/
```

- A module is reached through its barrel (`@/modules/rooms`), never its internals.
- `modules/` must not import from `app/`.
- `shared/` and `config/` must not import features.
- `app/app.tsx` is the composition root and the only file in `app/` that imports modules.

Inside a module: `api.ts`, `queries.ts` and `constants.ts` at the root, then `hooks/`,
`lib/` (pure, no React), `components/` and `pages/`.

**Go** is gofumpt-formatted with `golangci-lint` clean. The hub and the API are
separate modules on purpose; shared wire types live in `backend/pkg/wire`.

**Contracts** (`frontend/packages/contracts`) describe the HTTP wire in zod and are
shared by the web app and any future client, so they depend on zod and nothing else.

## Pull requests

Say what changed and why, and what a reviewer should look at first. If the change is
mostly moves or renames, say so — it saves the reviewer reading a large diff line by
line. Screenshots help for anything visual.

CI runs lint, format, vet, build, the Go tests, the module boundary check, and the web
lint and build. Everything should be green before review.
