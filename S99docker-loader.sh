#!/bin/sh

# Start all Docker Compose stacks in /usr/local/etc/init.d
# executing them in numerical order.

COMPOSE_BIN="/usr/local/bin/docker compose"
STACKS_DIR='/volmain/.@docker_compose/stacks'

start_stack() {
    stack_dir="$1"
    [ -d "$stack_dir" ] || return
    echo "Starting Docker Compose stack in '$stack_dir'..."
    $COMPOSE_BIN \
        --project-directory="$stack_dir" \
        up \
        --detach \
        --remove-orphans \
        --wait-timeout 30
}

stop_stack() {
    stack_dir="$1"
    [ -d "$stack_dir" ] || return
    echo "Stopping Docker Compose stack in '$stack_dir'..."
    $COMPOSE_BIN \
        --project-directory="$stack_dir" \
        down \
        --remove-orphans
}

case "$1" in
    start)
        for stack in "$STACKS_DIR"/*/; do
            start_stack "$stack"
        done
        ;;
    stop)
        for stack in "$STACKS_DIR"/*/; do
            stop_stack "$stack"
        done
        ;;
    reload|restart)
        for stack in "$STACKS_DIR"/*/; do
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
