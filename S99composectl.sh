#!/bin/sh

# Manage Docker Compose stacks via composectl CLI

COMPOSECTL="${COMPOSECTL:-/usr/local/bin/composectl}"

case "$1" in
    start|stop|restart|reload|down|list)
        exec "$COMPOSECTL" "$1"
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|reload}"
        exit 1
        ;;
esac
