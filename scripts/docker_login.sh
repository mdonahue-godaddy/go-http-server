#!/usr/bin/env bash
set -e -x -o pipefail

echo "BOF: ${0}"

#eval $(aws ecr get-login --no-include-email --region ${AWS_REGION})
#eval $(aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com)
aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

echo "EOF: ${0}"
