#!/usr/bin/env sh

# Get project root directory

ROOT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; cd .. ; pwd -P )";

# Run tests with coverage

COVERAGE_DIR="${ROOT_DIR}/coverage"

rm -rf "${COVERAGE_DIR}";

mkdir -p "${COVERAGE_DIR}";

go test -count=1 -race -coverpkg=./... -covermode=atomic -coverprofile="${COVERAGE_DIR}/coverage.cov" ./...;

if [ $? -ne 0 ]; then
    rm -rf "${COVERAGE_DIR}";

    exit 1;
fi

go tool cover -func="${COVERAGE_DIR}/coverage.cov";

rm -rf "${COVERAGE_DIR}";