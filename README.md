# container-wrapper

The goal of this project is to provide a simple shell interface and containerized environment for a C/C++ dev
environment.

## Goals 

- [ ] `start [CONTAINER]`: run container in background

- [ ] `shell [CONTAINER]`: start shell in container

- [ ] `setup [LOCATION]`: build/pull Dockerfile, setup run script

- [ ] `config [CONTAINER]`: manage multiple environments with single container wrapper
  - list mounted directories
  - list container source
  - (low priority) list container environments
  

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
