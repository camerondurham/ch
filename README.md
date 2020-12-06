# ch: container helper

The goal of this project is to provide a simple shell interface and containerized environment for a C/C++ dev
environment.

[![GitHub Super-Linter](https://github.com/camerondurham/ch/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)

## Goals

- [x] `create ENVIRONMENT`: create new environment config
  - [x] replace if it already exists
  - [x] build image or pull from Docker repository
    
- [x] `delete ENVIRONMENT`: delete environment

- [ ] `start ENVIRONMENT`: run container in background

- [ ] `shell ENVIRONMENT`: start shell in container

- [ ] `stop ENVIRONMENT`: shutdown running container

- [x] `list [CONTAINER]`: list container environment details, no args prints all details
  - [x] list mounted directories
  - [x] list container source

- [ ] unit test each command

## Spec

```shell script
# create environment
ch create ENVIRONMENT_NAME {--file DOCKERFILE|--image DOCKER_IMAGE} [--volume PATH_TO_DIRECTORY] [--shell SHELL_CMD]

ch create --file ./env/Dockerfile --shell /bin/bash --volume ./project/files/ --name cs104

# start container
# [docker run ]
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

support the following `docker build` flags:

```shell script
  -f, --file string             Name of the Dockerfile (Default is 'PATH/Dockerfile')
      --image string            Name of the Docker container to pull from Docker Hub
```

support the following `docker run` flags:

```shell script
      --cap-add list                   Add Linux capabilities
      --cap-drop list                  Drop Linux capabilities
      --name string                    Assign a name to the container
      --privileged                     Give extended privileges to this container
  -v, --volume list                    Bind mount a volume
  -w, --workdir string                 Working directory inside the container
```

### start

start docker container in background and save container id to config file

### shell

run docker shell in docker environment

### stop

stop running container/environment

### list

list all saved configs

### root

root command should allow you to set some defaults:

- `--set-default`: set default environment to start if you regularly use one


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
