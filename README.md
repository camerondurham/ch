# ch: container helper

A simple Docker interface to manage multiple containerized development environments. Provides a simple shell environment for separate development environments designed to use for C++ development in CSCI 104 but portable enough to use whichever Docker container you choose.

<a href="https://github.com/marketplace/actions/super-linter">
  <img align="left" src="https://github.com/camerondurham/ch/workflows/Lint%20Code%20Base/badge.svg" />
</a>

<a href="https://www.repostatus.org/#active">
  <img align="left" src="https://www.repostatus.org/badges/latest/active.svg" alt="Project Status: Active – The project has reached a stable, usable state and is being actively developed." />
</a>

<br>

## Commands

### create

create docker environment, specify Dockerfile to build or image to pull

**Supported Configuration:**

- ports
- bind mount volumes
- privileged
- security-opt

### start

start docker container in background and save container id to config file

### shell

run docker shell in docker environment

### stop

stop running container/environment

### list

list all saved configs

### running

list all running environments

## Examples

```shell script
# create environment
ch create ENVIRONMENT_NAME {--file DOCKERFILE|--image DOCKER_IMAGE} [--volume PATH_TO_DIRECTORY] [--shell SHELL_CMD] [--port HOST:CONTAINER] [--security-opt SECURITY_OPT]

ch create --image usccsci104/docker --shell /bin/bash --volume ./project/files/ --name cs104

ch create csci350 --image camerondurham/xv6-docker:latest --shell /bin/bash --security-opt seccomp:unconfined --port 7776:22 --port 7777:7777 --port 25000:25000 --cap-add SYS_PTRACE --privileged --replace

# start container - essentially docker run -d IMAGE 
ch start cs104

# get shell into environment - essentially docker exec -it CONTAINER_NAME
ch shell cs104

# stop container
ch stop cs104

# list all environments
ch list

# list all running environments
ch running
```
