#!/usr/bin/env sh

# Get project root directory

ROOT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; cd .. ; pwd -P )";

# Run tests

go test "${ROOT_DIR}"/... -race -count=1;