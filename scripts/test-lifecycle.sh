#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
just --justfile ./justfile setup
just --justfile ./justfile setup
output=$(printf '{"name":"Ada"}\n' | just --justfile ./justfile run)
test "$output" = "name: Ada"
