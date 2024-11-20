#!/bin/bash

IMAGE=mastrogiovanni/bridger:v0.0.1

docker run -it \
    -v $(pwd)/config.yaml:/app/config.yaml \
    -v /home/michele/.ssh/known_hosts:/root/.ssh/known_hosts \
    -v /home/michele/.ssh/id_rsa:/root/.ssh/id_rsa \
    -p 2345:2345 \
    -p 2346:2346 \
    -p 1234:1234 \
    $IMAGE bridger \
        /app/config.yaml \
        local:whoami:2345 \
        uat:calm:2346 \
        uat:db:1234
