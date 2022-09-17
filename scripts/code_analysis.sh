#!/usr/bin/env bash
set -e -x -o pipefail

echo "BOF: ${0}"

# Run code analysis reports.

source "$(dirname "$0")/common.sh"
: ${WORKSPACE:?}

cd $WORKSPACE/src/$REPONAME

pwd

echo "Starting go vet..."
go vet ./... 2> ${REPORTDIR}/govet-report.out || true

echo "Starting golint..."
golint ./... > ${REPORTDIR}/golint-report.out || true

echo "Starting golangci-lint..."
golangci-lint config --version --verbose
golangci-lint run --issues-exit-code 0 --out-format checkstyle > ${REPORTDIR}/checkstyle-result.xml

echo "Linters complete."

echo "EOF: ${0}"
