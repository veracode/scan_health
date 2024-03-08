#!/usr/bin/env bash

# Exit if any command fails
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[1;36m'
NC='\033[0m' # No Color


echo -e "${CYAN}Packaging for SAST scanning...${NC}"
go mod vendor
rm -f -- scan/veracode.zip
cd ..
zip -r scan_health/scan/veracode.zip scan_health -i "*.go" -i "**go.mod" -i "**go.sum"
cd scan_health


echo -e "\n${CYAN}Downloading Veracode CLI...${NC}"
cd scan
set +e # Ignore failure which happens if the CLI is the current latest version
curl -fsS https://tools.veracode.com/veracode-cli/install | sh
set -e
cd ..


echo -e "\n${CYAN}SAST Scanning with Veracode...${NC}"
./scan/veracode scan --type archive --source scan/veracode.zip --format table


echo -e "\n${CYAN}Container scanning with Veracode...${NC}"
set +e # Ignore failure
./scan/veracode scan --type image --source antfie/scan_health:latest --format table
set -e

echo -e "\n${CYAN}Container scanning with Scout...${NC}"
docker scout cves antfie/scan_health