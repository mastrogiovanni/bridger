#!/bin/bash

docker run -it \
    -v $(pwd)/config.yaml:/app/config.yaml \
    -v /home/michele/.ssh/known_hosts:/root/.ssh/known_hosts \
    -v /home/michele/.ssh/id_rsa:/root/.ssh/id_rsa \
    -p 2345:2345 \
    -p 1234:1234 \
    mastrogiovanni/bridger:main bridger \
        /app/config.yaml \
        uat:calm:2345 \
        uat:db:1234
