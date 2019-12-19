#!/usr/bin/env bash

set -x
set -e

rm -fr dist >/dev/null 2>&1 || true
mkdir -p dist

go build -v -o dist/castinator cmd/castinator/main.go
