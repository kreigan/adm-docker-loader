#!/bin/sh

set -e

DOCKER_COMPOSE_DIR='/volmain/.@docker_compose'

if [ -d "$DOCKER_COMPOSE_DIR" ]; then
    echo "Error: Directory '$DOCKER_COMPOSE_DIR' already exists."
    exit 1
fi

DOCKER_COMPOSE_BIN_DIR="$DOCKER_COMPOSE_DIR/bin"
mkdir -vp "$DOCKER_COMPOSE_BIN_DIR"

initScriptURL="https://raw.githubusercontent.com/kreigan/adm-docker-loader/main/S99docker-loader.sh"
loaderFile="$DOCKER_COMPOSE_BIN_DIR/start-stop.sh"
wget --tries=3 --output-document="$loaderFile" "$initScriptURL" 2>/dev/null
if [ $? -ne 0 ]; then
    echo "Error: Failed to download the Docker Compose loader script."
    exit 1
fi

DOCKER_COMPOSE_STACKS_DIR="$DOCKER_COMPOSE_DIR/stacks"
mkdir -vp "$DOCKER_COMPOSE_STACKS_DIR"

# Enable the Docker Compose loader script at system startup
chmod +x "$loaderFile"
initFile="/usr/local/etc/init.d/S99docker-loader.sh"
ln -s "$loaderFile" "$initFile"
