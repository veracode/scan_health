#!/usr/bin/env bash

# Exit if any command fails
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[1;36m'
NC='\033[0m' # No Color

# Download and extract the pipeline scan
if [ ! -f scan/pipeline-scan.jar ]; then
    echo -e "${CYAN}Downloading Veracode Pipeline Scanner...${NC}"
    curl -O https://downloads.veracode.com/securityscan/pipeline-scan-LATEST.zip
    unzip pipeline-scan-LATEST.zip pipeline-scan.jar
    mv pipeline-scan.jar scan/pipeline-scan.jar
    rm pipeline-scan-LATEST.zip
fi

echo -e "${CYAN}Packaging...${NC}"
go mod vendor
rm -f -- scan/veracode.zip
cd ..
zip -r scan_health/scan/veracode.zip scan_health -i "*.go" -i "**go.mod" -i "**go.sum"
cd scan_health

echo -e "${CYAN}SAST Scanning with Veracode...${NC}"
cd scan
java -jar pipeline-scan.jar --baseline_file baseline.json --file veracode.zip
cd ..