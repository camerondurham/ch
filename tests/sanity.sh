#!/bin/bash

TESTNAME="alpine-test1"

docker image rm -f "$TESTNAME" && docker system prune -f
touch ~/.ch.yaml

repo=$(dirname "$PWD")
DEBUG=1 go run "$repo/main.go" create "$TESTNAME" --file Dockerfile.alpine --shell /bin/sh
