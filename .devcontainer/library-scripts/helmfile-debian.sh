#!/usr/bin/env bash
#
# Syntax: ./helmfile-debian.sh [HELMFILE_VERSION]

HELMFILE_VERSION="v0.146.0"

set -e

if [ "$(id -u)" -ne 0 ]; then
    echo -e 'Script must be run as root. Use sudo, su, or add "USER root" to your Dockerfile before running this script'
    exit 1
fi

curl -Lo ./helmfile "https://github.com/helmfile/helmfile/releases/download/${HELMFILE_VERSION}/helmfile_linux_amd64"
chmod +x ./helmfile
mv ./helmfile /usr/local/bin/helmfile

# helmfile require helm diff plugin to work
helm plugin install https://github.com/databus23/helm-diff
