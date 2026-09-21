#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
EASYJSON="go run github.com/mailru/easyjson/easyjson@v0.9.0"

DTO_DIRS=(
	app/internal/campaign/presentation/http/dto
)

for dir in "${DTO_DIRS[@]}"; do
	for src in "$ROOT/$dir"/*.go; do
		[[ "$(basename "$src")" == *_easyjson.go ]] && continue
		(cd "$ROOT" && $EASYJSON -all "$dir/$(basename "$src")")
	done
done

echo "easyjson generated"
