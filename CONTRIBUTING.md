# Contributing to ch

## Development

### Adding New Commands

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

### Working with Go Modules

Go Module:

```shell script
# you don't have to run this since we already have a go.mod and go.sum file
go mod init github.com/<name>/<repo-name>

# add new library
go get <new dependency>

# organize modules and dependencies
go mod tidy

# remove dependency
go mod edit -dropreplace github.com/go-chi/chi
```

## Documentation

When writing instructions in the CLI and in the README, please follow syntax recommended by [google developer docs](https://developers.google.com/style/code-syntax)


Change package name:

```shell script
# change module name in all files
 find . -type f \( -name '*.go' -o -name '*.mod' \) -exec sed -i -e "s;container-helper;ch;g" {} +
```
