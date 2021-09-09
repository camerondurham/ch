# `ch` container helper

<div>
  <a href="github.com/camerondurham/ch">
    <img align="left" src="https://img.shields.io/github/v/release/camerondurham/ch?include_prereleases" />
  </a>

  <a href="github.com/camerondurham/ch">
    <img align="left" src="https://img.shields.io/github/go-mod/go-version/camerondurham/ch" />
  </a>

  <a href="github.com/camerondurham/ch">
    <img align="left" src="https://img.shields.io/github/workflow/status/camerondurham/ch/Build" />
  </a>

  <a href="https://github.com/marketplace/actions/super-linter">
    <img align="left" src="https://github.com/camerondurham/ch/workflows/Lint%20Code%20Base/badge.svg" />
  </a>

  <a href="https://github.com/camerondurham/ch">
    <img align="left" src="https://img.shields.io/github/downloads/camerondurham/ch/total" />
  </a>
</div>

<br>

`ch` is a command-line interface for using Docker containers as development environment and easily enter and exit
a container shell to compile and debug code. The tool provides a simple Docker interface to manage configuration
for multiple environments in a `~/.ch.yaml` (macOS/Linux) or `%userprofile%\.ch.yaml` (windows) file.
`ch` mimics the the `docker exec -it` command for its shell environment and supports the same but not all flags as
docker or docker-compose.

This was designed to generalize how we develop C++ code in CSCI 104 to be portable enough to use
whichever Docker container you choose. Of course, this project would not be possible without the reference
of [docker/cli](https://github.com/docker/cli) for examples of how to use the Docker Engine API.


## Installation

You can follow the installation instructions below to install
`ch` or see
the [csci104/docker](https://github.com/csci104/docker) repository
instructions and automated setup scripts for a C++ based
development environment.

### Prerequisites


Please make sure that your machine meets the requirements for Docker Desktop:

<a href="https://docs.docker.com/docker-for-windows/install/" target="_blank">Windows host:</a>

- Windows 10 64-bit: (Build 18362 or later)
  - WSL2 container backend

<a href="https://docs.docker.com/docker-for-mac/install/" target="_blank">Mac host:</a>

- Intel:
  - Mac hardware must be a 2010 or newer model
  - macOS must be version 10.13 or newer
  - 4 GB RAM minimum
- Apple Silicon (i.e. M1 chip):
  - Rosetta emulated terminal
    - for instructions on how to setup a Rosetta emulated terminal, see
      <a href="https://osxdaily.com/2020/11/18/how-run-homebrew-x86-terminal-apple-silicon-mac/" target="_blank">instructions here</a>
      to run Terminal through Rosetta.

### Step 0: Install WSL2 (Windows only)

If you are using macOS or Linux operating system, you can skip this section.
If you are running Windows, you must install the Windows Subsystem for Linux 2 (WSL2) before installing Docker.

Follow the instructions below to install WSL2 on your machine: <a href="https://docs.microsoft.com/windows/wsl/install-win10" target="_blank">Windows Subsystem for Linux Installation Guide</a>

### Step 1: Install Docker

Install Docker Desktop from <a href="https://www.docker.com/products/docker-desktop" target="_blank">the website</a>.

### Step 2: Install `ch`

Run the following commands below to download and run the install script for your operating system.

#### Windows

Run PowerShell as Admin and execute this command to download and run the install script for Windows:

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/camerondurham/ch/main/scripts/install-ch.ps1'))
```

You can check out the source code [here](https://github.com/camerondurham/ch/blob/main/scripts/install-ch.ps1).

You may need to restart your machine or log out so `ch` is added to your `Path`.

#### macOS/Linux

Run in your preferred Terminal to download and run the install script for Unix:

```bash
 bash <(curl -s https://raw.githubusercontent.com/camerondurham/ch/main/scripts/install-ch.sh)
 ```

 You can check out the source code [here](https://github.com/camerondurham/ch/blob/main/scripts/install-ch.sh).

Depending on your default shell (usually `bash` or `zsh`), you will have to source your `~/.bashrc` or `~/.zshrc` to add
`ch` to your `PATH`.

### Step 3: Setup Your First Environment

See [commands documentation](#commands) or the [example commands](#more-examples) for how to create your first
environment.

#### Create the CSCI104 Environment

Where `csci104-work` is your homework folder in the current directory.  Alternatively, you can provide the absolute path
to wherever your homework is on your machine.

This environment is based on this repository: [csci104/docker](https://github.com/csci104/docker)

```shell
# create the environment
ch create csci104 --image usccsci104/docker:20.04 --volume csci104-work:/work  --security-opt seccomp:unconfined --cap-add SYS_PTRACE --shell /bin/bash

# autostart and open a shell into the container
ch shell csci104 --force-start
```

#### Create the CSCI 350 Environment

Where `csci350-work` is your homework folder in the current directory. Alternatively, you can provide the absolute path
to wherever your homework is on your machine. For Windows, your volume command should look like `--volume "C:\Users\user\path\to\csci350:/work"`, on macOS your command should look like `--volume /Users/username/path/to/csci350:/work`.

This environment is based on the this repository: [camerondurham/cs350-docker](https://github.com/camerondurham/cs350-docker)

```shell
# create the environment
ch create csci350 --image camerondurham/cs350-docker:latest  --volume csci350-work:/xv6_docker --security-opt seccomp:unconfined --port 7776:22 --port 7777:7777 --port 25000:25000 --cap-add SYS_PTRACE --shell /bin/bash --privileged

# autostart and open a shell into the container
ch shell csci350 --force-start
```

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

### `ch update`

update an environment's image

If you specified a Dockerfile when using `ch create`, this will re-build that image.
If you specified a remote container registry when using `ch create`, it will try
to re-pull the image from that path.

```txt
Usage:
  ch update [ENVIRONMENT_NAME]
```

## More Examples

```shell script
# create environment
ch create ENVIRONMENT_NAME {--file DOCKERFILE|--image DOCKER_IMAGE} [--volume PATH_TO_DIRECTORY] [--shell SHELL_CMD] [--port HOST:CONTAINER] [--security-opt SECURITY_OPT]

# create an environment with a non-dockerhub image
ch create al2 --image public.ecr.aws/amazonlinux/amazonlinux:2 --shell /bin/bash --volume ./project/files

# create a csci104 docker image
ch create csci104 --image usccsci104/docker --shell /bin/bash --volume ./project/files/

# start container - essentially docker run -d IMAGE
ch start cs104

# get shell into environment - essentially docker exec -it CONTAINER_NAME
ch shell cs104

# stop container
ch stop cs104

# update docker image for environment
ch update csci104

# list all environments
ch list

# list all running environments
ch running
```
