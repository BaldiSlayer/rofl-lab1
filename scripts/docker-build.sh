#!/usr/bin/env bash

PACKAGE="github.com/BaldiSlayer/rofl-lab1"
COMMIT_HASH="$(git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')

LDFLAGS=(
  "-X '${PACKAGE}/internal/version.CommitHash=${COMMIT_HASH}'"
  "-X '${PACKAGE}/internal/version.BuildTime=${BUILD_TIMESTAMP}'"
)

go build -o /bin/backend -ldflags="${LDFLAGS[*]}" /app/cmd/backend/backend.go
