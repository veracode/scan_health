#!/usr/bin/env bash

# Exit if any command fails
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[1;36m'
NC='\033[0m' # No Color

mkdir -p scan
rm -f -- scan/veracode.zip


echo -e "${CYAN}Packaging for SAST scanning...${NC}"
go mod vendor
cd ..
zip -r scan_health/scan/veracode.zip scan_health -i "*.go" -i "**go.mod" -i "**go.sum"
cd scan_health


echo -e "\n${CYAN}Downloading the Veracode CLI...${NC}"
cd scan
set +e # Ignore failure which happens if the CLI is the current latest version
curl -fsS https://tools.veracode.com/veracode-cli/install | sh
set -e
cd ..


echo -e "\n${CYAN}SAST Scanning with Veracode...${NC}"
./scan/veracode static scan --baseline-file sast_baseline.json --results-file dist/sast_results.json scan/veracode.zip


echo -e "\n${CYAN}Container scanning with Veracode...${NC}"
./scan/veracode scan --type image --source antfie/scan_health:latest --output dist/container_results.json


echo -e "\n${CYAN}Container scanning with Scout...${NC}"
docker scout cves antfie/scan_health


echo -e "\n${CYAN}Generating SBOMs...${NC}"
./scan/veracode sbom --type archive --source scan/veracode.zip --output dist/src.sbom.json
./scan/veracode sbom --type image --source antfie/scan_health:latest --output dist/container.sbom.json