# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.203.0/containers/ubuntu/.devcontainer/base.Dockerfile

# [Choice] Ubuntu version (use hirsuite or bionic on local arm64/Apple Silicon): hirsute, focal, bionic
ARG VARIANT="hirsute"
FROM mcr.microsoft.com/vscode/devcontainers/base:0-${VARIANT}

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update \
    && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        build-essential libncurses-dev libssl-dev libffi-dev vim acl pip silversearcher-ag

# Install docker
RUN apt-get update && apt-get install -y --no-install-recommends \
    apt-transport-https ca-certificates curl gnupg2 gnupg-agent software-properties-common \
    && curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add - \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
    && apt update && apt -y install docker-ce docker-ce-cli containerd.io \
    && usermod -aG docker vscode \
    && newgrp docker

# Mount for docker-in-docker 
VOLUME [ "/var/lib/docker" ]

# Install docker compose
RUN curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose \
    && chmod +x /usr/local/bin/docker-compose