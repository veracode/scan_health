#!/usr/bin/env bash
set -e
./scan/veracode scan --type image --source antfie/scan_health:latest --output dist/container_results.json
