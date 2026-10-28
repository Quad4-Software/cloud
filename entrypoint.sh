#!/bin/sh
set -e

# Initialize a fresh config before starting the server so that a
# single `podman run` brings up a working instance.
if [ "$1" = "server" ] && [ ! -f /etc/quad4/quad4.yaml ]; then
    /usr/bin/quad4 init
fi

exec /usr/bin/quad4 "$@"
