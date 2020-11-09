#!/bin/bash

echo "Installing go tools"

LINT_VERSION=v1.31.0


ACTUAL_LINT_VERSION=$(golangci-lint version 2>&1)
if [[ $ACTUAL_LINT_VERSION =~ $LINT_VERSION ]]; then
    echo "golanci-lint is already installed"
else
    echo "updating golanci-lint..."
    GO111MODULE=on CGO_ENABLED=0 go get github.com/golangci/golangci-lint/cmd/golangci-lint@${LINT_VERSION}
fi