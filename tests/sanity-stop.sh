#!/bin/bash

TESTNAME="alpine-test$1"

docker image rm -f "$TESTNAME" && docker system prune -f
#touch ~/.ch.yaml

repo=$(dirname "$PWD")
DEBUG=1 go run "$repo/main.go" stop "$TESTNAME"
