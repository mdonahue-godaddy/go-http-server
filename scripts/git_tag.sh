#!/usr/bin/env bash
set -e -x -o pipefail

echo "BOF: ${0}"

if grep -q -i "alpine" /etc/os-release 2>/dev/null;
then
    alias awk=gawk
fi

function getVersion() {  
    if [[ ${BRANCH_NAME} == "master" ]]; then        
        echo "$(git semver next --dryrun)"
    else
        echo ${BRANCH_NAME}-${BUILD_NUMBER}
    fi
}

VERSION=$(getVersion)

echo "${VERSION}" > .git_tag
echo "Creating and pushing tag ${VERSION} for branch ${BRANCH_NAME}"

git tag ${VERSION}

if [[ ${BRANCH_NAME} == "master" ]]; then
    git push origin --tags
    #git push origin ${VERSION}
else
    echo "Skipping non master branch push"
fi

echo "EOF: ${0}"
