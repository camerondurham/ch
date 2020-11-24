# previously required:
# env GIT_TERMINAL_PROMPT=1

install:
	go get -t \
	github.com/camerondurham/container-wrapper \
	github.com/kelseyhightower/envconfig

dev-tools:
	go get -t \
	github.com/spf13/cobra/cobra
