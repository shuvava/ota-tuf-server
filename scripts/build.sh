#!/bin/sh

set -o errexit
set -o nounset

if [ -z "${OS:-}" ]; then
  echo "OS must be set"
  exit 1
fi
if [ -z "${ARCH:-}" ]; then
  echo "ARCH must be set"
  exit 1
fi
if [ -z "${VERSION:-}" ]; then
  echo "VERSION must be set"
  exit 1
fi
if [ -z "${COMMIT_HASH:-}" ]; then
  echo "COMMIT_HASH must be set"
  exit 1
fi

export CGO_ENABLED=0
export GOARCH="${ARCH}"
export GOOS="${OS}"
export GO111MODULE=on
export GOFLAGS="${GOFLAGS:-} -mod=${MOD}"

PACKAGE="$(go list -m)/pkg/version"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')


go install  \
  -installsuffix "static" \
  -ldflags "-X '${PACKAGE}.Version=${VERSION}' -X '${PACKAGE}.CommitHash=${COMMIT_HASH}' -X '${PACKAGE}.BuildDate=${BUILD_TIMESTAMP}'"  \
  "$@"
