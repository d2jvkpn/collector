#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

name=collector
mkdir -p configs logs data/mongo

export TAG=$1 APP_ENV=$2 PORT=$3
envsubst < ${_path}/docker_deploy.yaml > docker-compose.yaml

# docker-compose pull
[[ ! -z "$(docker ps --all --quiet --filter name=${name}_$APP_ENV)" ]] &&
  docker rm -f ${name}_$APP_ENV
# docker-compose down for running containers only, not stopped containers

docker-compose up -d
