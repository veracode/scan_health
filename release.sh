#!/usr/bin/env sh

VERSION="1.4"
FLAGS="-X main.AppVersion=$VERSION -s -w"

GOOS=darwin GOARCH=arm64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_health-mac-arm64" .
GOOS=darwin GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_health-mac-amd64" .
GOOS=linux GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_health-linux-amd64" .
GOOS=windows GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_health-win.exe" .

docker build -t antfie/scan_health:$VERSION .
docker build -t antfie/scan_health .
docker push antfie/scan_health:$VERSION
docker push antfie/scan_health