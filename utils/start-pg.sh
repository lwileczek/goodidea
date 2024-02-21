#!/bin/bash -xe

podman run \
    -v $(dirname $PWD)/initdb:/docker-entrypoint-initdb.d:ro,Z \
    --publish 5555:5432 \
    --detach \
    --name=pg \
    --env POSTGRES_USER=fairy \
    --env POSTGRES_PASSWORD=goodidea \
    --env POSTGRES_DB=tasks \
    docker.io/postgres:16-bookworm

# --mount type=bind,src=$(dirname $PWD)/initdb,destination=/docker-entrypoint-initdb.d,ro=true,bind-propagation=shared \
