#!/bin/bash
# https://github.com/golang/mock
go install github.com/golang/mock/mockgen@v1.6.0


cd ../cmd/util

mockgen -source ./docker_api.go -package=mocks -destination ./mocks/DockerAPI.go DockerAPI
mockgen -source ./docker.go -package=mocks -destination ./mocks/DockerClient.go DockerClient
mockgen -source ./validate.go -package=mocks -destination ./mocks/Validate.go Validate
