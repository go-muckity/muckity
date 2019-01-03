#!/usr/bin/env bash

set -e

source scripts/.functions.sh

if [ "$DEV_MONGO" != "docker" ]; then
    echo "Not gonna kill anything but a docker stack."; exit 1
fi

_cleanup
