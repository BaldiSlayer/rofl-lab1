#!/usr/bin/env bash

PACKAGE="github.com/BaldiSlayer/rofl-lab1"
VERSION="$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"
COMMIT_HASH="$(git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')

LDFLAGS=(
  "-X '${PACKAGE}/internal/version.Version=${VERSION}'"
  "-X '${PACKAGE}/internal/version.CommitHash=${COMMIT_HASH}'"
  "-X '${PACKAGE}/internal/version.BuildTime=${BUILD_TIMESTAMP}'"
)

go build -v -ldflags="${LDFLAGS[*]}" ./...
