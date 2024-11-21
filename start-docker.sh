#!/bin/bash

IMAGE=mastrogiovanni/bridger:v0.0.1

docker run -it \
    -v $(pwd)/config.yaml:/app/config.yaml \
    -v /home/user/.ssh/known_hosts:/root/.ssh/known_hosts \
    -v /home/user/.ssh/id_rsa:/root/.ssh/id_rsa \
    -p 9000:9000 \
    $IMAGE bridger \
        /app/config.yaml \
        dev:api-service:9000
