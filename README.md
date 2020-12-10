# ch: container helper [WIP]

A simple Docker interface to manage multiple containerized devleopment environments. Provides a simple shell environment for separate development environments designed to use for C++ development in CSCI 104 but portable enough to use whichever Docker container you choose.

<a href="https://github.com/marketplace/actions/super-linter">
  <img align="left" src="https://github.com/camerondurham/ch/workflows/Lint%20Code%20Base/badge.svg" />
</a>

<a href="https://www.repostatus.org/#wip">
  <img align="left" src="https://www.repostatus.org/badges/latest/wip.svg" alt="Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public."/>
</a>


</br>

## Goals

- [x] `create ENVIRONMENT` create new environment config
  - [x] replace if it already exists
  - [x] build image from Dockerfile
  - [x] pull image from Docker repository
  
- [x] `delete ENVIRONMENT` delete environment

- [x] `start ENVIRONMENT` run container in background

- [ ] `shell ENVIRONMENT` start shell in container
  - [ ] run commands in container
  - [ ] attach interactive terminal to container
  

- [x] `stop ENVIRONMENT` shutdown running container

- [x] `list [CONTAINER]` list container environment details, no args prints all details
  - [x] list mounted directories
  - [x] list container source

- [ ] unit test each command

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

## Development

Install cobra dependencies: (required to generate new commands)

```shell script
go get github.com/spf13/cobra/cobra
```

Add new cobra command

```shell script
# add new subcommand
cobra add <child command> -p <parent command>
cobra add childCommand -p 'parentCommand'
```

Add or adjust `~/.cobra.yaml` file for your name, license, year, etc. [Docs](https://github.com/spf13/cobra/blob/master/cobra/README.md)

Go Module:

```shell script
# already created
go mod init github.com/<name>/<repo-name>

# add new library
go get <new dependency>

# tidy your porject
go mod tidy

# remove dependency
go mod edit -dropreplace github.com/go-chi/chi
```

Write Documentation:

Follow syntax recommended by [google developer docs](https://developers.google.com/style/code-syntax)


Change package name:

```shell script
# change module name in all files
 find . -type f \( -name '*.go' -o -name '*.mod' \) -exec sed -i -e "s;container-helper;ch;g" {} +
```
