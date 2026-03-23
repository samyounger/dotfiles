#!/usr/bin/env bash

set -euo pipefail

repo_root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
if [ -z "$repo_root" ]; then
  printf 'Error: run this script from inside the git repository.\n' >&2
  exit 1
fi

cd "$repo_root"

if [ ! -d .githooks ]; then
  printf 'Error: missing %s/.githooks\n' "$repo_root" >&2
  exit 1
fi

if ! command -v gitleaks >/dev/null 2>&1; then
  printf 'Error: gitleaks is required. Install it through the Go CLI first, then re-run this script.\n' >&2
  exit 1
fi

chmod +x .githooks/pre-commit .githooks/pre-push
git config --local core.hooksPath .githooks

printf 'Enabled repo-managed git hooks in %s/.githooks\n' "$repo_root"
