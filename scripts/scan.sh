#!/usr/bin/env bash

ESCAPE=$'\e'

# Download and extract the pipeline scan
if [ ! -f scan/pipeline-scan.jar ]; then

    echo "${ESCAPE}[0;36mDownloading Veracode Pipeline Scanner...${ESCAPE}[0m"
    curl -O https://downloads.veracode.com/securityscan/pipeline-scan-LATEST.zip
    unzip pipeline-scan-LATEST.zip pipeline-scan.jar
    mv pipeline-scan.jar scan/pipeline-scan.jar
    rm pipeline-scan-LATEST.zip
fi

echo "${ESCAPE}[0;36mPackaging...${ESCAPE}[0m"
go mod vendor
rm scan/veracode.zip
cd ..
zip -r scan_health/scan/veracode.zip scan_health -x "scan_health/dist/*" -x "scan_health/cache/*" -x "scan_health/scan/*" -x "scan_health/scripts/*" -x "scan_health/docs/*" -x "scan_health/.git/*" -x "scan_health/.idea/*" -x "**/.DS_Store" -x "scan_health/.*"
cd scan_health

echo "${ESCAPE}[0;36mScanning...${ESCAPE}[0m"
cd scan
java -jar pipeline-scan.jar --baseline_file baseline.json --file veracode.zip
cd ..