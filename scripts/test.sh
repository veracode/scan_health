#!/usr/bin/env bash

# Exit if any command fails
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[1;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}Linting...${NC}"
gofmt -s -w .

echo -e "${CYAN}Running go vet...${NC}"
go vet ./...

echo -e "${CYAN}Running gosec...${NC}"
set +e # Don't fail on error for this
gosec ./...
set -e

echo -e "${CYAN}Testing...${NC}"

go test -v -coverprofile coverage.out -failfast -shuffle on ./...
go tool cover -html=coverage.out -o coverage.html