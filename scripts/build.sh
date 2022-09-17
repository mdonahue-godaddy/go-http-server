#!/usr/bin/env bash
set -e -x -o pipefail

echo "BOF: ${0}"

source "$(dirname "$0")/common.sh"
: ${WORKSPACE:?}

cd $WORKSPACE/src/$REPONAME

pwd

#VERSION=${1:-undefined}
#echo "Building version ${VERSION} ..."

GOROOT=/usr/local/go make clean tag
if [ "${1}" == "prerelease" ]; then
    GOROOT=/usr/local/go goreleaser release --snapshot --debug --rm-dist
elif [ "${1}" == "release" ]; then
    GOROOT=/usr/local/go goreleaser release --debug --rm-dist
else
    echo "INVALID PARAMETER: '${1}'"
fi

echo "EOF: ${0}"