#!/usr/bin/env bash
set -e -x -o pipefail

echo "BOF: ${0}"

# Set up common environment variables.

export REPONAME=github.secureserver.net/urlshortener/forwarding-svc

# Ensure that GOPATH is set. If it is initially unset or empty, use the value of WORKSPACE
# (which will always be set in Jenkins environments). If WORKSPACE is not set, fall back to
# the default value for Unix systems.
export GOPATH=${GOPATH:-${WORKSPACE:-$HOME/go}}

# Ensure that WORKSPACE is set. If it is initially unset or empty, use the first element of
# the colon-separated GOPATH as the value.
export WORKSPACE=${WORKSPACE:-${GOPATH%%:*}}

# Ensure that the local Go binaries are available. This will have no ill effect
# if the WORKSPACE binary path is already part of PATH.
export PATH=$PATH:${WORKSPACE}/bin

# If GO_HOME is set, then we are probably on Jenkins and need to include it in the path
if [ ! -z "$GO_HOME" ]; then
  export PATH=$PATH:$GO_HOME/bin
fi

export REPODIR=${WORKSPACE}/src/${REPONAME}
export REPORTDIR=${WORKSPACE}/reports

if [ ! -d "${REPORTDIR}" ]; then
  mkdir -p "${REPORTDIR}"
fi

echo "EOF: ${0}"
