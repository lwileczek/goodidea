#!/bin/bash -xe

sudo docker run \
    -v $(dirname $PWD)/initdb:/docker-entrypoint-initdb.d \
    --publish 5555:5432 \
    --detach \
    --name=pg \
    --env POSTGRES_USER=fairy \
    --env POSTGRES_PASSWORD=goodidea \
    --env POSTGRES_DB=tasks \
    postgres:16-bookworm
