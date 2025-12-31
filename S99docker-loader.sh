#!/bin/sh

# Start all Docker Compose stacks in /usr/local/etc/init.d
# executing them in numerical order.

COMPOSE_BIN="/usr/local/bin/docker compose"
DOCKER_LOADER_DIR='/volmain/.@docker_compose'
CONFIG_FILE="$DOCKER_LOADER_DIR/config.conf"
ENV_FILE="$DOCKER_LOADER_DIR/.env"
STACKS_DIR="$DOCKER_LOADER_DIR/stacks"

COMPOSE_COMMON_ARGS="--env-file $ENV_FILE"

COMPOSE_UP_ARGS="--detach \
--wait-timeout 30 \
--pull always"

if [ -f "$CONFIG_FILE" ]; then
    echo "Loading configuration from '$CONFIG_FILE'..."
    source "$CONFIG_FILE" 2>/dev/null
fi

getStackName() {
    basename "$1" | cut -d'-' -f2-
}

composeStack() {
    local stackDir="$1"
    local action="$2"
    
    [ -d "$stackDir" ] || return
    
    local stackName=$(getStackName "$stackDir")
    echo "Executing 'docker compose $action' for stack '$stackName' in '$stackDir'..."
    
    local actionArgs=""
    case "$action" in
        up)
            actionArgs="$COMPOSE_UP_ARGS"
            ;;
        down)
            actionArgs="$COMPOSE_DOWN_ARGS"
            ;;
        *)
            echo "Error: Unknown action '$action'"
            return 1
            ;;
    esac

    if [ -f "$stackDir/$action-args" ]; then
        echo "Overriding run arguments from stack-specific $action-args file..."
        actionArgs=$(cat "$stackDir/$action-args")
    fi

    $COMPOSE_BIN \
        $COMPOSE_COMMON_ARGS \
        --project-directory="$stackDir" \
        --project-name "$stackName" \
        "$action" \
        $actionArgs
}

start_stack() {
    composeStack "$1" up
}

stop_stack() {
    composeStack "$1" down
}

getStacks() {
    find "$STACKS_DIR" -maxdepth 1 -mindepth 1 -type d -name '??-*' | sort
}

case "$1" in
    start)
        for stack in $(getStacks); do
            start_stack "$stack"
        done
        ;;
    stop)
        for stack in $(getStacks); do
            stop_stack "$stack"
        done
        ;;
    reload|restart)
        for stack in $(getStacks); do
            stop_stack "$stack"
            start_stack "$stack"
        done
        ;;
    *)
        echo "Usage: $0 {start|stop|reload|restart}"
        exit 1
        ;;
esac

exit 0
