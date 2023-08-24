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
zip -r scan_health/scan/veracode.zip scan_health -x "scan_health/dist/*" -x "scan_health/cache/*" \
  -x "scan_health/scan/*" -x "scan_health/scripts/*" -x "scan_health/docs/*" -x "scan_health/.git/*" \
  -x "scan_health/.idea/*" -x "**/.DS_Store" -x "scan_health/.*"
cd scan_health

echo -e "${CYAN}Scanning...${NC}"
cd scan
java -jar pipeline-scan.jar --baseline_file baseline.json --file veracode.zip
cd ..