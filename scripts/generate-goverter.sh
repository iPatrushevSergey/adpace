#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

CONVERTER_PKGS=(
	./app/internal/campaign/adapter/repository/postgres/converter/...
)

failed=0
for pkg in "${CONVERTER_PKGS[@]}"; do
	if ! (cd "$ROOT" && go generate "$pkg"); then
		echo "goverter failed: $pkg" >&2
		failed=1
	fi
done

if [[ "$failed" -ne 0 ]]; then
	echo "goverter finished with errors (install: go install github.com/jmattheis/goverter/cmd/goverter@latest)" >&2
	exit 1
fi

echo "goverter generated"
