#!/bin/bash

set -euo pipefail

TOP_DIR="$( cd "$(dirname "$0")" ; pwd -P )"
cd "${TOP_DIR}"

VERSION_TAG=${VERSION_TAG:-"local-$(whoami)"}
GIT_HASH=$(git rev-parse HEAD)
BUILD_TIME=$(TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')
BUILD_NUMBER=${BUILD_NUMBER:-0}
VERSION_INFO="{\"build_number\":${BUILD_NUMBER},\"version\":\"${VERSION_TAG}\",\"built_time\":\"${BUILD_TIME}\",\"git_hash\":\"${GIT_HASH}\"}"
echo "${VERSION_INFO}"
echo "${VERSION_INFO}" >release.json

BUILD_MODE=${BUILD_MODE:-"release"}

GOOS=linux GOARCH=amd64 go build -ldflags "\
  -s -w \
  -X main.Version=${VERSION_TAG} \
  -X main.GitHash=${GIT_HASH} \
  -X main.BuildTime=${BUILD_TIME} \
  -X main.BuildMode=${BUILD_MODE}" \
  -o mbtilesConverter_linux

GOOS=windows GOARCH=amd64 go build -ldflags "\
  -s -w \
  -X main.Version=${VERSION_TAG} \
  -X main.GitHash=${GIT_HASH} \
  -X main.BuildTime=${BUILD_TIME} \
  -X main.BuildMode=${BUILD_MODE}" \
  -o mbtilesConverter_windows.exe
