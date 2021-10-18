# CHANGELOG

## v0.3.4 (unreleased)

**NOTE**:
This changelog is for upcoming release `v0.3.4`. Anyone who builds `ch` from source has these updates but the binary has not been released yet.

- `ch update` command
  - will re-build or re-pull the latest image for the environment
- improved error when install script is run on Linux
  - install script relies on `unzip` command on Unix which is installed by default on macOS
  - if `unzip` command fails, script prints out how to install on Debian-based and Arch-based distros
- dependency upgrades for `containerd` security vulnerabilities
- upgrade from go `1.16` to `1.17`
