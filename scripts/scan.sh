#!/usr/bin/env bash

# Exit if any command fails
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[1;36m'
NC='\033[0m' # No Color

zipFilePath=".veracode_scan/packages/veracode-auto-pack-scan_health-go.zip"

rm -rf -- .veracode_scan
mkdir .veracode_scan


echo -e "\n${CYAN}Downloading the Veracode CLI...${NC}"
cd .veracode_scan
set +e # Ignore failure which happens if the CLI is the current latest version
curl -fsS https://tools.veracode.com/veracode-cli/install | sh
set -e
cd ..


echo -e "\n${CYAN}Packaging for SAST scanning...${NC}"
./.veracode_scan/veracode package --trust --source . --output .veracode_scan/packages


echo -e "\n${CYAN}SAST Scanning with Veracode (Pipeline)...${NC}"
./scan/veracode static scan --baseline-file sast_baseline.json \
                            --results-file dist/sast_results.json $zipFilePath


echo -e "\n${CYAN}Container scanning with Veracode...${NC}"
./scan/veracode scan --type image --source antfie/scan_health:latest \
                     --output dist/container_results.json


echo -e "\n${CYAN}Container scanning with Scout...${NC}"
docker scout cves antfie/scan_health


echo -e "\n${CYAN}Generating SBOMs...${NC}"
./scan/veracode sbom --type archive --source $zipFilePath --output dist/src.sbom.json
./scan/veracode sbom --type image --source antfie/scan_health:latest --output dist/container.sbom.json


echo -e "\n${CYAN}SAST Scanning with Veracode (Policy)...${NC}"
# Refer to this page: https://central.sonatype.com/artifact/com.veracode.vosp.api.wrappers/vosp-api-wrappers-java/versions
apiWrapperVersion="24.4.13.0"
curl "https://repo1.maven.org/maven2/com/veracode/vosp/api/wrappers/vosp-api-wrappers-java/$apiWrapperVersion/vosp-api-wrappers-java-$apiWrapperVersion.jar" --output scan/veracode-api.jar
java -jar scan/veracode-api.jar -action uploadandscan \
     -appname "scan_health" -createprofile false \
     -version `date "+%Y-%m-%d %H:%M:%S"` -filepath .veracode_scan/packages