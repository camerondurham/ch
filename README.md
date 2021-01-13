# ch: container helper

**This project is in a beta stage. Please report any issues and I'll be happy to address!**

A simple Docker interface to manage multiple containerized development environments. Provides a simple shell environment for separate development environments designed to use for C++ development in CSCI 104 but portable enough to use whichever Docker container you choose.

<a href="https://github.com/marketplace/actions/super-linter">
  <img align="left" src="https://github.com/camerondurham/ch/workflows/Lint%20Code%20Base/badge.svg" />
</a>

<a href="https://www.repostatus.org/#wip">
  <img align="left" src="https://www.repostatus.org/badges/latest/wip.svg" alt="Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public."/>
</a>


<br>

## Status

- [x] `create ENVIRONMENT` create new environment config
  - [x] replace if it already exists
  - [x] build image from Dockerfile
  - [x] pull image from Docker repository
  - [x] support volume mounts and path checking on Windows and macOS
  
- [x] `delete ENVIRONMENT` delete environment

- [x] `start ENVIRONMENT` run container in background

- [x] `shell ENVIRONMENT` start shell in container
  - [x] run commands in container
  - [x] attach interactive terminal to container
  - [x] attach volumes to containers
    - [x] MVP on Windows
    - [x] MVP on macOS/Unix
  - [x] add options for debugging in containers
    - [x] security options (`seccomp:unconfined`)
    - [x] add capacities (`SYS_PTRACE`)
  

- [x] `stop ENVIRONMENT` shutdown running container

- [x] `list [CONTAINER]` list container environment details, no args prints all details
  - [x] list mounted directories
  - [x] list container source

- [ ] unit tests
  - [x] `create`
  - [ ] `delete`
  - [ ] `start`
  - [ ] `shell`
  - [ ] `stop`
  - [ ] `list`

## Spec

```shell script
# create environment
ch create ENVIRONMENT_NAME {--file DOCKERFILE|--image DOCKER_IMAGE} [--volume PATH_TO_DIRECTORY] [--shell SHELL_CMD]

ch create --file ./env/Dockerfile --shell /bin/bash --volume ./project/files/ --name cs104

# start container
# [docker run]
ch start cs104

# get shell into environment
#  [docker exec]
ch shell cs104

# stop container
ch stop cs104

# list all environments
ch list

# set default container when you don't provide any args
ch --set-default cs104
```

### create

create docker environment, specify Dockerfile to build or image to pull

### start

start docker container in background and save container id to config file

### shell

run docker shell in docker environment

### stop

stop running container/environment

### list

list all saved configs
