#!/usr/bin/env bash
# End-to-end smoke test for btd: build the binary, scaffold a project with
# kumi, apply the tomo/hachigo includes, and verify the result — all in a
# throwaway directory.
#
# Prerequisites: go, git, node (>=20). kumi is installed by `btd new` if
# missing; kumi's own prerequisites (templ, air, goose) must be present.

set -euo pipefail

cd "$(dirname "$0")/.."

BINDIR="$(mktemp -d)"
WORK="$(mktemp -d)"
trap 'rm -rf "$BINDIR" "$WORK"' EXIT

echo "==> building btd"
go build -o "$BINDIR/btd" ./cmd/btd

echo "==> btd new kumi smokeapp --include tomo,hachigo"
(cd "$WORK" && "$BINDIR/btd" new kumi smokeapp --include tomo,hachigo)

echo "==> verifying scaffold"
test -d "$WORK/smokeapp"
test -f "$WORK/smokeapp/.btd/manifest.json"
grep -q '^\.tomo/$' "$WORK/smokeapp/.gitignore"

echo "==> btd install tomo (inside the project)"
(cd "$WORK/smokeapp" && "$BINDIR/btd" install tomo)

echo "==> btd list"
"$BINDIR/btd" list

echo "smoke test passed"
