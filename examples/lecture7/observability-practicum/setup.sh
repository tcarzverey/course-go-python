#!/usr/bin/env bash
# Run once after cloning to generate go.sum files for both services.
set -euo pipefail

BASE="$(cd "$(dirname "$0")" && pwd)"

echo "==> go mod tidy: url-shortener"
cd "$BASE/url-shortener"
go mod tidy

echo "==> go mod tidy: stats-service"
cd "$BASE/stats-service"
go mod tidy

echo ""
echo "✅  Done. You can now run: docker compose up --build -d"
