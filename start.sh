#!/bin/bash

cat env.sh

# export variables defined in env.sh
set -a && source env.sh && set +a
docker stack deploy -c docker-compose.yml dr-gateway