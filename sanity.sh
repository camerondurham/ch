#!/bin/bash

TESTNAME="alpine-test1"

docker image rm -f "$TESTNAME" && docker system prune -f
touch ~/.ch.yaml

go run main.go create "$TESTNAME" --file tests/Dockerfile.alpine --shell /bin/sh 
