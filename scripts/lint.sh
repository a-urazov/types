#!/bin/bash

# Script to run golangci-lint with proper error handling

set -euo pipefail

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "golangci-lint could not be found. Installing..."
    # Install golangci-lint
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.60.3
fi

echo "Running golangci-lint..."
golangci-lint run ./... --timeout=5m

echo "Linting completed successfully!"