# !/usr/bin/env sh

ESCAPE=$'\e'
export VERSION="1.14"

./build.sh && \

echo "${ESCAPE}[0;32mReleasing v${VERSION}...${ESCAPE}[0m" && \

docker build -t antfie/scan_health:$VERSION . && \
docker build -t antfie/scan_health . && \
docker push antfie/scan_health:$VERSION && \
docker push antfie/scan_health && \

echo "${ESCAPE}[0;32mRelease Success${ESCAPE}[0m"