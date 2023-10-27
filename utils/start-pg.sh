#!/bin/bash -xe

sudo docker run \
    -v $(dirname $PWD)/initdb:/docker-entrypoint-initdb.d \
    --publish 5555:5432 \
    --detach \
    --name=pg \
    --env POSTGRES_USER=cody \
    --env POSTGRES_PASSWORD=sucks \
    --env POSTGRES_DB=tasks \
    postgres:16-bookworm
