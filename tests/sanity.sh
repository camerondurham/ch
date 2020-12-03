#!/bin/bash

TESTNAME="alpine-test1"

docker image rm -f "$TESTNAME" && docker system prune -f
touch ~/.ch.yaml

# TODO: update main.go path
go run main.go create "$TESTNAME" --file tests/Dockerfile.alpine --shell /bin/sh
