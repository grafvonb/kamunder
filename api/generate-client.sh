#!/bin/bash
# ./scripts/swagger-openapi-goclient/4-generate-client.sh ./api/uamoim-local/openapi-app.json ./internal/clients/uamoim/apiapp/client.gen.go apiapp
# ./scripts/swagger-openapi-goclient/4-generate-client.sh ./api/uamoim-local/openapi-imxlogin.json ./internal/clients/uamoim/apiimxlogin/client.gen.go apiimxlogin

set -euo pipefail

need() { command -v "$1" >/dev/null 2>&1 || { echo "missing tool: $1" >&2; exit 127; }; }
need oapi-codegen

src="${1:-${SRC:-}}"
out="${2:-${OUT:-}}"
pkg="${3:-${PKG:-}}"

oapi-codegen -generate types,client -package "$pkg" -o "$out" "$src"