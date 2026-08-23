#!/usr/bin/env bash
set -euo pipefail

tag="${1:?usage: release-notes.sh <tag> [previous-tag]}"
previous="${2:-}"

if [ -z "$previous" ]; then
  previous="$(git describe --tags --abbrev=0 --exclude="$tag" "$tag^" 2>/dev/null || true)"
fi

if [ -n "$previous" ]; then
  range="$previous..$tag"
else
  range="$tag"
fi

subjects="$(git log --no-merges --reverse --pretty=format:'%s' "$range")"

section() {
  local heading="$1" pattern="$2"
  local body
  body="$(printf '%s\n' "$subjects" | grep -E "$pattern" | sed -E 's/^[a-z]+(\([^)]*\))?!?: *//' | sed 's/^/- /' || true)"
  if [ -n "$body" ]; then
    printf '### %s\n\n%s\n\n' "$heading" "$body"
  fi
}

{
  section "Breaking changes" '^[a-z]+(\([^)]*\))?!:'
  section "Features" '^feat(\([^)]*\))?: '
  section "Fixes" '^fix(\([^)]*\))?: '
  section "Other" '^(chore|docs|refactor|perf|test|build|ci|style)(\([^)]*\))?: '

  printf '### Hub install\n\n'
  printf 'Raspberry Pi (64-bit):\n\n'
  printf '```sh\n'
  printf 'curl -fsSL -o avahub https://github.com/%s/releases/download/%s/avahub-linux-arm64\n' "${GITHUB_REPOSITORY:-xenomech/ava}" "$tag"
  printf 'chmod +x avahub\n'
  printf '```\n\n'
  printf 'Verify against `checksums.txt` before running it. See `docs/releasing.md` for the full install.\n'

  if [ -n "$previous" ]; then
    printf '\n**Full changelog**: https://github.com/%s/compare/%s...%s\n' "${GITHUB_REPOSITORY:-xenomech/ava}" "$previous" "$tag"
  fi
}
