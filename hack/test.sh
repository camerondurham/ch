#!/bin/bash
set -x

REPOSITORY=$PWD

if [ $(basename $REPOSITORY) != "ch" ]; then
  echo "must run script from ch repository root"
  exit 1
fi

# this script is a WIP goal is to run tests

docker build -t ch-test-img:v0 -f $REPOSITORY/hack/Dockerfile $REPOSITORY

