# ch: container helper

A simple Docker interface to manage multiple containerized development environments. Provides a simple shell environment for separate development environments designed to use for C++ development in CSCI 104 but portable enough to use whichever Docker container you choose.
<div>
  <a href="github.com/camerondurham/ch">
    <img align="left" src="https://img.shields.io/github/v/release/camerondurham/ch?include_prereleases" />
  </a>

  <a href="github.com/camerondurham/ch">
    <img align="left" src="https://img.shields.io/github/go-mod/go-version/camerondurham/ch" />
  </a>

  <a href="https://github.com/marketplace/actions/super-linter">
    <img align="left" src="https://github.com/camerondurham/ch/workflows/Lint%20Code%20Base/badge.svg" />
  </a>

  <a href="https://www.repostatus.org/#active">
    <img align="left" src="https://www.repostatus.org/badges/latest/active.svg" alt="Project Status: Active – The project has reached a stable, usable state and is being actively developed." />
  </a>
  
  <a href="https://github.com/camerondurham/ch">
    <img align="left" src="https://img.shields.io/github/downloads/camerondurham/ch/total" />
  </a>
</div>
<br>

## What is this?

What's the use case for this tool? Good question! This tool is designed to make it easier to use a specific, isolated development environment. For classes
such as CSCI 104 and CSCI 350 at USC, the legacy way of writing code in the class was using a large VM image inside Virtual Box,
or if you're lucky, VMWare. A more efficient and arguably smoother workflow involves setting using a Docker container with the class's compilers and
development tools installed. `ch` offers a consistent interface to configure and access these environments. See below for the commands to create
environments for these classes. All you have to do is run the command and the tool will download the required dependencies from DockerHub.

### Create the CSCI104 Environment

Where `csci104-work` is your homework folder in the current directory. This environment
is based on this repository: [csci104/docker](https://github.com/csci104/docker)

```shell
# create the environment
ch create cs104 \
  --image usccsci104/docker \
  --volume ./csci104-work:/work \
  --shell /bin/bash
  
# autostart and open a shell into the container
ch shell csci104 --force-start
```

### Create the CSCI 350 Environment

Where `csci350-work` is your homework folder in the current directory. This environment
is based on the this repository: [camerondurham/cs350-docker](https://github.com/camerondurham/cs350-docker)

```shell
# create the environment
ch create csci350 \
  --image camerondurham/cs350-docker:latest \
  --volume ./cs350-work:/xv6_docker \
  --security-opt seccomp:unconfined \
  --port 7776:22 \
  --port 7777:7777 \
  --port 25000:25000 \
  --cap-add SYS_PTRACE \
  --shell /bin/bash \
  --privileged

# autostart and open a shell into the container
ch shell csci350 --force-start
```

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

## More Examples

```shell script
# create environment
ch create ENVIRONMENT_NAME {--file DOCKERFILE|--image DOCKER_IMAGE} [--volume PATH_TO_DIRECTORY] [--shell SHELL_CMD] [--port HOST:CONTAINER] [--security-opt SECURITY_OPT]

ch create csci104 --image usccsci104/docker --shell /bin/bash --volume ./project/files/

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

