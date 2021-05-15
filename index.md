# `ch` container helper

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

  <a href="https://github.com/camerondurham/ch">
    <img align="left" src="https://img.shields.io/github/downloads/camerondurham/ch/total" />
  </a>
</div>

<br>

`ch` is a command-line interface for using Docker containers as development environment. The tool provides a simple
Docker interface to manage multiple containerized development environments. Like the `docker exec -it`, the CLI has a
shell environment. This was designed to generalize how we develop C++ code in CSCI 104 to be portable enough to use
whichever Docker container you choose. Of course, this project would not be possible without the reference
of [docker/cli](https://github.com/docker/cli) which is how I learned how to use the Docker Engine API.

## What is this?

What's the use case for this tool? Good question! This tool is designed to make it easier to use a specific, isolated development environment. For classes
such as [CSCI 104](https://bytes.usc.edu/cs104/) (Data Structures and Algorithms) and CSCI 350 (Operating Systems) at USC, the legacy way of writing code in the class was using a large VM image inside Virtual Box,
or if you're lucky, VMWare. A more efficient and arguably smoother workflow involves setting using a Docker container with the class's compilers and
development tools installed. `ch` offers a consistent interface to configure and access these environments. See below for the commands to create
environments for these classes.


## Commands

### `ch create`

create docker environment, specify Dockerfile to build or image to pull

```txt
Usage:
  ch create ENVIRONMENT_NAME {--file DOCKERFILE|--image DOCKER_IMAGE} [--volume PATH_TO_DIRECTORY] [--shell SHELL_CMD] [[--cap-add cap1] ...] [[--security-opt secopt1] ...] [flags]

Flags:
      --cap-add stringArray        special capacity to add to Docker Container (syscalls)
      --context string             context to build Dockerfile (default ".")
  -f, --file string                path to Dockerfile, should be relative to context flag
  -h, --help                       help for create
  -i, --image string               image name to pull from DockerHub
  -p, --port stringArray           bind host port(s) to container
      --privileged                 run container as privileged (full root/admin access)
      --replace                    replace environment if it already exists
      --security-opt stringArray   security options
      --shell string               default shell to use when logging into environment (default "/bin/sh")
      --version                    version for create
  -v, --volume stringArray         volume to mount to the working directory

```

### `ch start`

start docker container in background and save container id to config file

```txt
Usage:
  ch start ENVIRONMENT_NAME [flags]
Flags:
  -h, --help      help for start
  -v, --version   version for start
```

### `ch shell`

run docker shell in docker environment

```txt
Usage:
  ch shell ENVIRONMENT_NAME [flags]

Flags:
  -f, --force-start   autostart the environment if not running
  -h, --help          help for shell
  -v, --version       version for shell
```

### `ch stop`

stop running container/environment

```txt
Usage:
  ch stop ENVIRONMENT_NAME [flags]

Flags:
  -h, --help      help for stop
  -v, --version   version for stop
```

### `ch list`

list all saved configs

```txt
Usage:
  ch list [ENVIRONMENT_NAME] [flags]

Flags:
  -h, --help      help for list
  -v, --version   version for list
```

### `ch running`

list all running environments

```txt
Usage:
  ch running
```

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


<!-- ### Jekyll Themes

Your Pages site will use the layout and styles from the Jekyll theme you have selected in your [repository settings](https://github.com/camerondurham/ch/settings). The name of this theme is saved in the Jekyll `_config.yml` configuration file.

### Support or Contact

Having trouble with Pages? Check out our [documentation](https://docs.github.com/categories/github-pages-basics/) or [contact support](https://support.github.com/contact) and we’ll help you sort it out. -->
