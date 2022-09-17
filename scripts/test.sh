#!/usr/bin/env bash
set -e -x -o pipefail

echo "BOF: ${0}"

# Run the test and coverage reports. May be run locally and is used by the Jenkinsfile.

source "$(dirname "$0")/common.sh"
: ${WORKSPACE:?}

cd ${WORKSPACE}

pwd

echo "Updating go-junit-report ..."
#go get -u github.com/jstemmer/go-junit-report
go install github.com/jstemmer/go-junit-report@latest

echo "Updating gocover-cobertura ..."
#go get -u github.com/t-yuki/gocover-cobertura
go install github.com/t-yuki/gocover-cobertura@latest

cd ${REPODIR}

pwd

echo "Running tests ..."
go test -v ./... 2>&1 | ${WORKSPACE}/bin/go-junit-report > "${REPORTDIR}/junit_results.xml"

echo "Generating coverage report in ${REPORTDIR}..."
go test -coverprofile="${REPORTDIR}/cover.out" ./...

${WORKSPACE}/bin/gocover-cobertura < "${REPORTDIR}/cover.out" > "${REPORTDIR}/cobertura-coverage.xml"

echo "EOF: ${0}"
