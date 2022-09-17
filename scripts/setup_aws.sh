#!/usr/bin/env bash
set -e -x -o pipefail

echo "BOF: ${0}"

export SRCDIR=${1}

# ensure that pip is available
function ensurePipExists()
{
    if ! [[ -x "$(command -v pip)" ]]; then
        curl -s https://bootstrap.pypa.io/get-pip.py -o get-pip.py
        python get-pip.py --user
        rm get-pip.py
    fi
}

# check default python version
PYVERS=$(python -c 'import platform; major, minor, patch = platform.python_version_tuple(); print(major)')

# create venv with highest available python version
if [[ $PYVERS -eq 2 ]]; then
    # see if they have python 3 available
    if [[ -x "$(command -v python3)" ]]; then
        python3 -m venv .venv
    else
        # if they only have python 2, locally install virtualenv
        ensurePipExists
        #pip install -q virtualenv --user
        pip install virtualenv --user
        virtualenv .venv
    fi
elif [[ $PYVERS -eq 3 ]]; then
    python -m venv .venv
fi

# activate venv
source .venv/bin/activate

# install AWS CLI, OKTA Processor, and Sceptre
ensurePipExists

pip install --upgrade pip
pip install wheel
pip install -q -r ${SRCDIR}/requirements.txt
#pip install -q .

# exit venv and add okta alias to its activation script
deactivate

echo "alias okta='OKTA_DOMAIN="godaddy.okta.com"; KEY=\$(openssl rand -hex 18); eval \$(aws-okta-processor authenticate -e -o \$OKTA_DOMAIN -u \$USER -k \$KEY)'" >> .venv/bin/activate

# venv should be ready to use :)
echo VirtualEnv created, to activate run \"source .venv/bin/activate\"

echo "EOF: ${0}"
