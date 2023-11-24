#!/usr/bin/env sh

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )/../" &> /dev/null && pwd )

docker-compose -f ${SCRIPT_DIR}/docker-compose.yaml stop;

docker-compose -f ${SCRIPT_DIR}/docker-compose.yaml rm -f --all;