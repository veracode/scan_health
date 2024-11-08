#!/usr/bin/env bash

# Exit if any command fails
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[1;36m'
NC='\033[0m' # No Color

rm -rf dist
mkdir dist
./scripts/test.sh

if [[ -z "${VERSION}" ]]; then
    VERSION="0.0"
fi

FLAGS="-X main.AppVersion=$VERSION -s -w"

echo -e "\n${CYAN}Building v${VERSION}...${NC}"
GOOS=darwin GOARCH=arm64 go build -ldflags="$FLAGS" -buildvcs=false -trimpath -o "dist/scan_health-mac-arm64" .
GOOS=darwin GOARCH=amd64 go build -ldflags="$FLAGS" -buildvcs=false -trimpath -o "dist/scan_health-mac-amd64" .
GOOS=linux GOARCH=amd64 go build -ldflags="$FLAGS" -buildvcs=false -trimpath -o "dist/scan_health-linux-amd64" .
GOOS=windows GOARCH=amd64 go build -ldflags="$FLAGS" -buildvcs=false -trimpath -o "dist/scan_health.exe" .

echo -e "\n${GREEN}Build Success${NC}"