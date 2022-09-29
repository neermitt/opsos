# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.245.2/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.19, 1.18, 1-bullseye, 1.19-bullseye, 1.18-bullseye, 1-buster, 1.19-buster, 1.18-buster
ARG VARIANT="bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/base:${VARIANT}

ARG KIND_VERSION="v0.15.0"
ARG YQ_VERSION="v4.2.0"
ARG HELMFILE_VERSION="v0.144.0"


RUN apt-get update \
     && apt-get install -y --no-install-recommends make xz-utils\
     && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

COPY library-scripts/*.sh /tmp/library-scripts/

RUN /tmp/library-scripts/kind-debian.sh "${KIND_VERSION}" \
     && /tmp/library-scripts/helmfile-debian.sh "${HELMFILE_VERSION}" \
     && /tmp/library-scripts/yq-debian.sh "${YQ_VERSION}" \
     && rm -rf /tmp/library-scripts/
