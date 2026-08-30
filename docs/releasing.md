# Releasing

For maintainers. Two things ship on different clocks: the server and web app deploy
from a branch, and the hub binaries come from a tag.

## Promoting to production

`stage` is preproduction, `main` is production. Work lands on `stage` first, gets
verified there, and is then promoted as a whole:

```bash
gh pr create --base main --head stage --title "Promote stage"
```

Merging that is the production deploy, so check preproduction is healthy first. The
two branches should never diverge — if `stage` is not an ancestor of `main`,
something was merged to the wrong place.

## Cutting a hub release

The hub runs on someone else's Raspberry Pi, so it ships as a downloadable binary
rather than a deploy. Tag from `main`:

```bash
git checkout main && git pull
git tag v0.1.0
git push origin v0.1.0
```

Any `v*` tag runs `.github/workflows/release.yml`:

1. **verify** — `make test` and `make check-boundary` across all three Go modules
2. **hub** — cross-compiles `avahub` and `avactl` for `linux/arm64`, `linux/amd64`
   and `darwin/arm64`, with the version stamped in
3. **publish** — SHA256 checksums, notes generated from commit subjects, GitHub Release

Nothing publishes if `verify` fails. Confirm the release lists `avahub-linux-arm64`
and `checksums.txt` — those are what `scripts/install-hub.sh` downloads.

Versions are `vMAJOR.MINOR.PATCH`. Until the wire protocol between hub and API is
stable, treat any change to `backend/pkg/wire` as a minor bump and expect hubs to
need updating alongside the server.

### Release notes

`.github/scripts/release-notes.sh` groups commit subjects by conventional-commit
prefix into Breaking / Features / Fixes / Other. Preview them:

```bash
.github/scripts/release-notes.sh v0.1.0
```

With no earlier tag it walks the whole history, so pass a starting point to trim it:

```bash
.github/scripts/release-notes.sh v0.1.0 v0.0.9
```

## Deployment configuration

Environment variables, broker authentication and the SPA cache rules are documented
in [self-hosting.md](self-hosting.md) — the deployment needs the same settings anyone
self-hosting does.
