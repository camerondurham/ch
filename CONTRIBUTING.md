# Contributing to ch

## Development

### Development Environment

It is highly recommended to use GoLand from JetBrains for development. If you are a student,
you can use GoLand for free with an educational licence.

1. Install go locally or let GoLand install for you
2. Sync dependencies (this can likely be done if you simply `go build`, or simply let GoLand install for you)
3. Build, test, code, and repeat

On a Linux/macOS machine, you should be able to use `make` command to build, test, and package the repository.


### Adding New Commands

Install cobra dependencies: (required to generate new commands)

```bash
go install github.com/spf13/cobra/cobra@v1.3.0
```

Add new cobra command

```bash
# add new subcommand
cobra add <child command> -p <parent command>
cobra add childCommand -p 'parentCommand'
```

Add or adjust `~/.cobra.yaml` file for your name, license, year, etc. [Docs](https://github.com/spf13/cobra/blob/master/cobra/README.md)

### Working with Go Modules

Go Module:

```bash
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

```bash
# change module name in all files
 find . -type f \( -name '*.go' -o -name '*.mod' \) -exec sed -i -e "s;container-helper;ch;g" {} +
```


## Releasing New Versions

> Note: This is mainly a reminder for me since
> I can't remember how on Earth to git tag stuff sometimes!


```bash
# releasing a new version
git tag -a v1.4 -m "Version notes here"
```

To delete tagging mistakes locally:

```bash
# delete local tag
git tag -d tagname
```

To delete remote tags:

```bash
# delete a tag already pushed to github
git push --delete origin tagname
```
