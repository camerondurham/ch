# previously required:
# env GIT_TERMINAL_PROMPT=1

build:
	go build

cleanup:
	go mod tidy

get-tools:
	go get -t \
	github.com/spf13/cobra/cobra
