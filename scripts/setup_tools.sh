#!/usr/bin/env bash
set -e -x -o pipefail

echo "BOF: ${0}"

export SRCDIR=${1}

echo "Installing golangci-lint"
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint --version

echo "Installing gocov"
go install github.com/axw/gocov/gocov@latest

echo "Installing gocov-xml"
go install github.com/AlekSi/gocov-xml@latest

echo "Installing go2xunit"
go install github.com/tebeka/go2xunit@latest
go2xunit -version

echo "Installing golint"
go install golang.org/x/lint/golint@latest

echo "Installing mockgen"
go install github.com/golang/mock/mockgen@latest
mockgen -version

if (git semver &> /dev/null); then
    echo "Already Installed: git-semver"
else 
    echo "Installing git-semver"
    mkdir -p ~/tools 
    mkdir -p ~/.local/bin
    cd ~/tools
    git clone https://github.com/markchalloner/git-semver.git || true
    git-semver/install.sh
    cd -
fi

if [ -f ~/.local/bin/envsubst ]; then
    echo "Already Installed: envsubst"
else 
    echo "Installing envsubst"
    mkdir -p ~/.local/bin
    curl -L https://github.com/a8m/envsubst/releases/download/v1.1.0/envsubst-`uname -s`-`uname -m` -o envsubst
    chmod +x envsubst
    mv envsubst ~/.local/bin
fi

echo "Go Lang version"
go version

echo "Docker version"
docker version

echo "Make version"
make --version

echo "AWS version"
aws --version

echo "EOF: ${0}"
