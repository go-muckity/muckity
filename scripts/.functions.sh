#!/usr/bin/env bash

export DEV_MONGO=${DEV_MONGO:-$1}
export COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME:-muckity}
export MONGODB_VERSION=${MONGODB_VERSION:-4.0}
export MONGO_EXPRESS_VERSION=${MONGO_EXPRESS_VERSION:-0.49}

_cleanup() {
    docker rmi -f muckity:local-dev || echo "no muckity:local-dev to remove"
    if [ "$DEV_MONGO" = "docker" ]; then
        if [ "$TRAVIS_BUILD_NUMBER" != "" ]; then
            COMPOSE_PROJECT_NAME=muckity_travis
        fi
        docker-compose -f docker/dev-stack.yml down
    fi
}

_make_mongodb() {
    if [ "$DEV_MONGO" != "docker" ]; then
        echo "function called out of context: ${DEV_MONGO}!"; exit 1;
    fi
    if [ "$TRAVIS_BUILD_NUMBER" != "" ]; then
        COMPOSE_PROJECT_NAME=muckity_travis
    fi
    docker-compose -f docker/dev-stack.yml up $@
}
