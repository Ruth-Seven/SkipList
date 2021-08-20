#!/usr/bin/env bash

set -e

if ! type golangci-lint &>/dev/null; then
	echo "golangci-lint is not found, please install it"
	exit 1
fi

golangci-lint run -c .//.golangci-lint.yaml . --timeout 5m
