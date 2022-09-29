#!/usr/bin/env bash
#
# Syntax: ./yq-debian.sh [YQ_VERSION]

YQ_VERSION="v4.2.0"

set -e

BINARY="yq_linux_amd64"

if [ "$(id -u)" -ne 0 ]; then
    echo -e 'Script must be run as root. Use sudo, su, or add "USER root" to your Dockerfile before running this script'
    exit 1
fi

curl -Ls yq.tar.gz "https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/${BINARY}.tar.gz" | tar xz
mv "${BINARY}" /usr/local/bin/yq
