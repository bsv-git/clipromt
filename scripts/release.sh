#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
version=0.1.0
mkdir -p dist
for arch in amd64 arm64; do
    stage=$(mktemp -d)
    CGO_ENABLED=0 GOOS=linux GOARCH="$arch" go build -buildvcs=false -trimpath -ldflags="-s -w -X main.version=$version" -o "$stage/clipromt" ./cmd/clipromt
    cp README.md LICENSE "$stage/"
    tar -czf "dist/clipromt_${version}_linux_${arch}.tar.gz" -C "$stage" clipromt README.md LICENSE
    rm -rf "$stage"
done
(cd dist && sha256sum clipromt_0.1.0_linux_*.tar.gz > SHA256SUMS)
