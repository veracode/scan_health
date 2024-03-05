#!/usr/bin/env bash

# Exit if any command fails
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[1;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}Updating dependencies...${NC}"
go get -t -u ./...
go mod tidy
go mod vendor