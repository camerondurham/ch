# CHANGELOG

## v0.3.7

- Add `--platform` flag to specify which platform to pull when creating environment (thanks to @a-harhar!)

## v0.3.6

- Dependabot and README changes

## v0.3.5

- remove "latest version" checks from other commands (`ch list` and `ch running`)
- add "latest version" check logic to `ch upgrade` command

## v0.3.4

- `ch update` command
  - will re-build or re-pull the latest image for the environment
- `ch upgrade` command
  - unfortunately, this won't upgrade `ch` for you but will at least tell you the OS-specific command to run in your terminal
- improved error when install script is run on Linux
  - install script relies on `unzip` command on Unix which is installed by default on macOS
  - if `unzip` command fails, script prints out how to install on Debian-based and Arch-based distros
- dependency upgrades for `containerd` security vulnerabilities
- upgrade from go `1.16` to `1.17`
