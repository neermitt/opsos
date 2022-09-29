#!/usr/bin/env bash
#
# Syntax: ./kind-debian.sh [KIND_VERSION]

KIND_VERSION="v0.15.0"

set -e

if [ "$(id -u)" -ne 0 ]; then
    echo -e 'Script must be run as root. Use sudo, su, or add "USER root" to your Dockerfile before running this script'
    exit 1
fi

curl -Lo ./kind "https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-linux-amd64"
chmod +x ./kind
mv ./kind /usr/local/bin/kind
