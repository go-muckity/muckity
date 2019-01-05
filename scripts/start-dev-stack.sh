#!/usr/bin/env bash

set -e

source scripts/.functions.sh

if [ "$DEV_MONGO" = "" ]; then
    DEV_MONGO="server"
fi

_cleanup
_make_mongodb -d
