#!/usr/bin/env sh

# make sh play nicely
set -euo

# This only works on Linux. It is only expected to be used as part of CI/CD in a docker container.
echo "9e83cc74e67a945ad770c1b9851d9c048ecd327bcd5971c852ffaa53329cf69c  build/godownloader_0.1.0_Linux_x86_64.tar.gz" > godownloader.checksums

mkdir -p build
wget -P build https://github.com/goreleaser/godownloader/releases/download/v0.1.0/godownloader_0.1.0_Linux_x86_64.tar.gz
sha256sum -c godownloader.checksums && tar -xzf build/godownloader_0.1.0_Linux_x86_64.tar.gz -C build
