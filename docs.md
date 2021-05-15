# Documentation

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


