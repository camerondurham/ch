# Getting Started

## Installation

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

```bash
# create the environment
ch create csci104 --image usccsci104/docker:20.04 --volume csci104-work:/work  --security-opt seccomp:unconfined --cap-add SYS_PTRACE --shell /bin/bash

# autostart and open a shell into the container
ch shell csci104 --force-start
```

#### Create the CSCI 350 Environment

Where `csci350-work` is your homework folder in the current directory. Alternatively, you can provide the absolute path
to wherever your homework is on your machine. For Windows, your volume command should look like `--volume "C:\Users\user\path\to\csci350:/work"`, on macOS your command should look like `--volume /Users/username/path/to/csci350:/work`.

This environment is based on the this repository: [camerondurham/cs350-docker](https://github.com/camerondurham/cs350-docker)

```bash
# create the environment
ch create csci350 --image camerondurham/cs350-docker:latest  --volume csci350-work:/xv6_docker --security-opt seccomp:unconfined --port 7776:22 --port 7777:7777 --port 25000:25000 --cap-add SYS_PTRACE --shell /bin/bash --privileged

# autostart and open a shell into the container
ch shell csci350 --force-start
```

## More Examples

```bash

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
