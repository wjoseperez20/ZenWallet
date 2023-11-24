#!/usr/bin/env sh

# Get project root directory

ROOT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 && cd .. && pwd -P )";

command -v golangci-lint && {
  golangci-lint run;

  exit $?;
}

# Validate Docker is installed

command -v docker || {
  echo " - [ ERROR ] Docker is NOT installed! Install it following the instructions in: https://docs.docker.com/desktop/mac/install/";

  exit 1;
}

# Run lint

docker run --rm -v "${ROOT_DIR}:/app:ro" --workdir="/app" golangci/golangci-lint:v1.52 golangci-lint run;