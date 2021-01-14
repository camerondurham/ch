#!/bin/bash

TESTNAME="alpine-test$1"

docker image rm -f "$TESTNAME" && docker system prune -f

DEBUG=1 go run ../../main.go start "$TESTNAME"
