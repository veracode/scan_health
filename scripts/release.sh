#!/usr/bin/env bash

# Exit if any command fails
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[1;36m'
NC='\033[0m' # No Color

if [ $# -eq 0 ]
  then
    echo -e "${RED}Error: No version number specified. Try something like \"release.sh x.y\".${NC}"
    exit
fi

export VERSION=$1

./scripts/build.sh

echo -e "${CYAN}Releasing v${VERSION}...${NC}"

docker pull alpine
docker build -t antfie/scan_health .
docker push antfie/scan_health

docker scout cves antfie/scan_health

echo -e "${GREEN}Release Success${NC}"